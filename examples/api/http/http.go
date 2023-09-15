package http

import (
	"context"
	"fmt"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

// Adding requestID in the fault context for logging purpose
func DecorateRequestMetadata(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rid := middleware.GetReqID(ctx)
		ctx = fctx.WithMeta(ctx, "request_id", rid)
		ctx = fctx.WithMeta(ctx, "http_method", r.Method)
		ctx = fctx.WithMeta(ctx, "request_path", r.URL.Path)
		ctx = fctx.WithMeta(ctx, "remote_ip", r.RemoteAddr)
		ctx = fctx.WithMeta(ctx, "protocol", r.Proto)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Adding any path variable in the fault context for logging purpose
func PathVariableAsFCtx(pathVarName, fctxName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if value := chi.URLParam(r, pathVarName); value != "" {
				r = r.WithContext(fctx.WithMeta(r.Context(), fctxName, value))
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func fctxToSlog(ctx context.Context) []any {
	faultFields := fctx.GetMeta(ctx)
	fields := make([]any, 0, len(faultFields))
	for k, v := range faultFields {
		fields = append(fields, slog.String(k, v))
	}
	return fields
}

func LoggerRequest(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				fields := fctxToSlog(r.Context())
				fields = append(fields, slog.Int("status", ww.Status()))
				fields = append(fields, slog.Duration("latency", time.Since(t1)))
				logger.Info("API Request", fields...)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}

}

func RespondWithError(
	logger *slog.Logger,
	err error,
	w http.ResponseWriter,
	r *http.Request,
) {
	tag := ftag.Get(err)

	attrs := fctxToSlog(r.Context())
	errStr := fmt.Sprintf("%+v", err)
	attrs = append(attrs, slog.String("error", err.Error()))
	logger.Error("\n"+errStr, attrs...)

	// Using tags to determine http status based on the error
	if tag == ftag.NotFound {
		http.Error(w, fmsg.GetIssue(err), 404)
		return
	}

	// Default to internal server error
	http.Error(w, http.StatusText(500), 500)
}
