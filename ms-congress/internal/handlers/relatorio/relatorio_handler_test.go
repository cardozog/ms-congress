package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestListarEventosRejectsInvalidFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		header     string
		query      string
		wantStatus int
	}{
		{name: "missing organizer id", wantStatus: http.StatusBadRequest},
		{name: "invalid start date", header: "1", query: "?dataDe=2026-13-01", wantStatus: http.StatusBadRequest},
		{name: "start date after end date", header: "1", query: "?dataDe=2026-11-20&dataAte=2026-11-19", wantStatus: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			context.Request = httptest.NewRequest(http.MethodGet, "/relatorios/eventos"+test.query, nil)
			if test.header != "" {
				context.Request.Header.Set("X-Organizador-ID", test.header)
			}

			NewHandler(nil).ListarEventos(context)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
		})
	}
}
