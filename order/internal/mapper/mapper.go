package mapper

import (
	"order_service/internal/infra"
	pb "order_service/proto/order_service"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPbOrder(order *infra.Order) *pb.Order {
	return &pb.Order{
		OrderId:       order.OrderID.String(),
		UserId:        order.UserID.String(),
		OrderAddress:  order.OrderAddress,
		OrderLocation: order.OrderLocation,
		OrderDate:     timestamppb.New(order.OrderDate),
		OrderTimeGap:  durationpb.New(order.OrderTimeGap),
		OrderStatus:   order.OrderStatus,
	}
}

func ToInfraOrder(order *pb.Order) *infra.Order {
	return &infra.Order{
		OrderID:       uuid.MustParse(order.OrderId), //FIXME: change to parse
		UserID:        uuid.MustParse(order.UserId),  //FIXME: change to parse
		OrderAddress:  order.OrderAddress,
		OrderLocation: order.OrderLocation,
		OrderDate:     order.OrderDate.AsTime(),
		OrderTimeGap:  order.OrderTimeGap.AsDuration(),
		OrderStatus:   order.OrderStatus,
	}
}

func ToPbOrders(orders []*infra.Order) []*pb.Order {
	pbOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = ToPbOrder(order)
	}
	return pbOrders
}

func ToInfraOrders(pbOrders []*pb.Order) []*infra.Order {
	infraOrders := make([]*infra.Order, len(pbOrders))
	for i, pbOrder := range pbOrders {
		infraOrders[i] = ToInfraOrder(pbOrder)
	}
	return infraOrders
}
