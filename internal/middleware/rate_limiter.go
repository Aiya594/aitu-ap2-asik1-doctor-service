package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func RateLimiterInterceptor(redisClient *redis.Client, logger *slog.Logger) grpc.UnaryServerInterceptor {
	rpmLimit := 100
	if v := os.Getenv("RATE_LIMIT_RPM"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			rpmLimit = n
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// If Redis is unavailable, skip rate limiting gracefully
		if redisClient == nil {
			return handler(ctx, req)
		}

		clientIP := extractClientIP(ctx)
		key := fmt.Sprintf("ratelimit:doctor:%s", clientIP)

		count, err := incrementSlidingWindow(ctx, redisClient, key, time.Minute)
		if err != nil {
			// Best-effort: log and allow through
			logger.Warn("rate limiter redis error, allowing request", "error", err, "client_ip", clientIP)
			return handler(ctx, req)
		}

		if count > int64(rpmLimit) {
			logger.Warn("rate limit exceeded",
				"client_ip", clientIP,
				"count", count,
				"limit", rpmLimit,
			)
			return nil, status.Errorf(codes.ResourceExhausted,
				"rate limit exceeded: %d requests per minute allowed; retry after %d seconds",
				rpmLimit, 60)
		}

		return handler(ctx, req)
	}
}

// incrementSlidingWindow uses INCR + EXPIRE to implement a fixed 1-minute window counter.
func incrementSlidingWindow(ctx context.Context, client *redis.Client, key string, window time.Duration) (int64, error) {
	pipe := client.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incrCmd.Val(), nil
}

func extractClientIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "unknown"
	}
	return p.Addr.String()
}
