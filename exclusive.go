package exclusive

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
)

type Configure struct {
	Responses map[int]any
}

type Doc struct {
	ID           string
	Tags         string
	Description  string
	IsDeprecated bool
	Request      any
	Response     any
}

func Swag(engine *gin.Engine, c *gin.Context, options *Configure) string {

	reflector := openapi3.Reflector{}

	for _, route := range engine.Routes() {

		if strings.Contains(route.Handler, "StaticFile") {
			continue
		}

		if strings.Contains(route.Path, c.Request.URL.String()) {
			continue
		}

		route.HandlerFunc(c)

		operation, _ := reflector.NewOperationContext(route.Method, route.Path)

		document, docExist := c.Get("doc")

		if docExist {

			operation.SetID(reflect.ValueOf(document).FieldByName("ID").String())

			operation.SetTags(reflect.ValueOf(document).FieldByName("Tags").String())
			operation.SetDescription(reflect.ValueOf(document).FieldByName("Description").String())
			operation.SetIsDeprecated(reflect.ValueOf(document).FieldByName("IsDeprecated").Bool())
			operation.AddReqStructure(reflect.ValueOf(document).FieldByName("Request").Interface())
			operation.AddRespStructure(reflect.ValueOf(document).FieldByName("Response").Interface())
		}

		for status, response := range options.Responses {
			operation.AddRespStructure(response, openapi.WithHTTPStatus(status))
		}

		err := reflector.AddOperation(operation)

		if err != nil {
			panic(err)
		}
	}

	schema, err := reflector.Spec.MarshalJSON()

	if err != nil {
		panic(err)
	}

	fmt.Println("Exclusive: Swagger documantation generated.")

	return string(schema)
}
