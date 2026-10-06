package merchantdocumenthandler

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_document"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-pkg/upload_image"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant_documents"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantDocument struct {
	Client      *grpc.ClientConn
	E           *echo.Echo
	Logger      logger.LoggerInterface
	UploadImage upload_image.ImageUploads
}

func RegisterMerchantDocumentHandler(deps *DepsMerchantDocument) {
	mapper := apimapper.NewMerchantDocumentResponseMapper()

	NewMerchantDocumentQueryHandleApi(&merchantDocumentQueryHandleDeps{
		client: pb_merchant_document.NewMerchantDocumentQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
	})

	NewMerchantDocumentCommandHandleApi(&merchantDocumentCommandHandleDeps{
		client:       pb_merchant_document.NewMerchantDocumentCommandServiceClient(deps.Client),
		router:       deps.E,
		logger:       deps.Logger,
		mapper:       mapper.CommandMapper(),
		upload_image: deps.UploadImage,
	})
}
