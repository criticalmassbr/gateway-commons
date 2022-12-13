package queue

import (
	"context"
	"encoding/json"

	utils "github.com/criticalmassbr/ms-utils"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Result struct {
	Message *amqp.Delivery
	Error   error
}

type QueueService interface {
	Publish(ctx context.Context, queueName string, publishing amqp.Publishing) error
	Rpc(ctx context.Context, queueName string, publishing amqp.Publishing) <-chan Result
}

func preparePublishing(c *fiber.Ctx, p *amqp.Publishing, opName string) {
	p.ContentType = "application/json"

	if p.Headers != nil {
		p.Headers["ClientSlug"] = c.Get("Slug")
	} else {
		p.Headers = amqp.Table{"ClientSlug": c.Get("Slug")}
	}

	p.Type = opName
}

func Publish(c *fiber.Ctx, q QueueService, p amqp.Publishing, queue, opName string) error {
	ctx, span := utils.Tracer.NewSpan(c.UserContext(), "fiber_utils", "Publish", oteltrace.WithAttributes(attribute.String("queue", queue)))
	defer span.End()

	preparePublishing(c, &p, opName)

	return q.Publish(ctx, queue, p)
}

func Rpc(c *fiber.Ctx, q QueueService, p amqp.Publishing, queue, opName string) (*amqp.Delivery, error) {
	_, span := utils.Tracer.NewSpan(c.UserContext(), "fiber_utils", "RpcService", oteltrace.WithAttributes(attribute.String("queue", queue)))
	defer span.End()

	preparePublishing(c, &p, opName)

	result := <-q.Rpc(c.UserContext(), queue, p)
	return result.Message, result.Error
}

func RpcToFiberResponse[T any](c *fiber.Ctx, q QueueService, p amqp.Publishing, queue, opName string) error {
	response, err := Rpc(c, q, p, queue, opName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	var res T

	if err = json.Unmarshal(response.Body, &res); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": res})
}
