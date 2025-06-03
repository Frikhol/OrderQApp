package entrypoint

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"notification_service/internal/config"
	impl "notification_service/internal/impl"
	"notification_service/internal/infra/broker"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var clients = sync.Map{}

func addClient(userID uuid.UUID, conn *websocket.Conn) {
	value, _ := clients.LoadOrStore(userID, []*websocket.Conn{})
	conns := value.([]*websocket.Conn)
	conns = append(conns, conn)
	clients.Store(userID, conns)
}

func removeClient(userID uuid.UUID, conn *websocket.Conn) {
	value, ok := clients.Load(userID)
	if !ok {
		return
	}
	conns := value.([]*websocket.Conn)
	var newConns []*websocket.Conn
	for _, c := range conns {
		if c != conn {
			newConns = append(newConns, c)
		}
	}
	if len(newConns) == 0 {
		clients.Delete(userID)
	} else {
		clients.Store(userID, newConns)
	}
}

func SendNotification(userID uuid.UUID, msg any) {
	value, ok := clients.Load(userID)
	if !ok {
		return
	}

	conns := value.([]*websocket.Conn)
	for _, conn := range conns {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Println("Write error:", err)
		}
	}
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(r.URL.Query().Get("user_id"))
	if userID == uuid.Nil {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	// 3. Upgrade до WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // на проде проверь origin!
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// 4. Сохраняем соединение
	addClient(userID, conn)
	log.Println("Client connected:", userID)

	// 5. Слушаем клиент → если отключился — удаляем
	go func() {
		defer func() {
			removeClient(userID, conn)
			conn.Close()
			log.Println("Client disconnected:", userID)
		}()
	}()
}

type Order struct {
	OrderID       uuid.UUID     `json:"order_id"`
	UserID        uuid.UUID     `json:"user_id"`
	AgentID       uuid.UUID     `json:"agent_id"`
	OrderAddress  string        `json:"order_address"`
	OrderLocation string        `json:"order_location"`
	OrderDate     time.Time     `json:"order_date"`
	OrderTimeGap  time.Duration `json:"order_time_gap"`
	OrderStatus   string        `json:"order_status"`
}

func Run(cfg *config.Config, logger *zap.Logger) error {

	broker, err := broker.New(logger, &cfg.RabbitMQ)
	if err != nil {
		logger.Fatal("failed to create broker", zap.Error(err))
	}
	defer broker.Close()

	// grpcServer := grpc.NewServer()
	// proto.RegisterNotificationServiceServer(grpcServer, handlers.New(impl.New(logger, broker)))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	// 	if err != nil {
	// 		logger.Fatal("failed to listen", zap.Error(err))
	// 	}
	// 	logger.Info("Notification service started", zap.String("port", cfg.GRPCPort))
	// 	if err := grpcServer.Serve(lis); err != nil {
	// 		logger.Error("failed to serve", zap.Error(err))
	// 	}
	// }()

	http.HandleFunc("/ws", WebSocketHandler)
	go func() {
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			logger.Error("failed to start http server", zap.Error(err))
		}
	}()

	go func() {
		msgs, err := broker.GetChannel().Consume(
			"queue_order_cancelled",
			"",
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			logger.Error("failed to consume messages", zap.Error(err))
		}

		for msg := range msgs {
			logger.Info("received cancelled order", zap.String("message", string(msg.Body)))
			order := Order{}
			err := json.Unmarshal(msg.Body, &order)
			if err != nil {
				logger.Error("failed to unmarshal message", zap.Error(err))
			}
			SendNotification(order.UserID, msg.Body)
		}
	}()

	go func() {
		err := impl.New(logger, broker).HandleOrderCancelledMessages()
		if err != nil {
			logger.Error("failed to handle messages", zap.Error(err))
		}
	}()

	go func() {
		err := impl.New(logger, broker).HandleOrderCompletedMessages()
		if err != nil {
			logger.Error("failed to handle messages", zap.Error(err))
		}
	}()

	<-done
	logger.Info("Notification service stopped")
	// grpcServer.Stop()

	return nil
}
