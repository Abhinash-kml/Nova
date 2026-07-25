package inventory

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, ic *ItemsController, inc *InventoryController) {
	// itemsPublic := router.Group("/public/items")
	// {

	// }

	itemsPrivate := router.Group("/private/items")
	{
		itemsPrivate.GET("", ic.GetAll)
		itemsPrivate.GET("/:id", ic.Get)
		itemsPrivate.POST("", ic.Create)
		itemsPrivate.PATCH("/:id", ic.Modify)
		itemsPrivate.DELETE("/:id", ic.Delete)
	}

	inventoryPublic := router.Group("/public/inventory")
	{
		inventoryPublic.GET("", inc.GetInventoryOfUser)
		inventoryPublic.DELETE("", inc.DeleteInventoryOfUser)
	}

	inventoryPrivate := router.Group("/private/inventory")
	{
		inventoryPrivate.GET("/:itemid", inc.GetInventoryItemOfUser)
		inventoryPrivate.POST("", inc.AddInventoryItemOfUser)
		inventoryPrivate.PATCH("/:itemid", inc.UpdateInventoryItemOfUser)
		inventoryPrivate.DELETE("/:itemid", inc.DeleteInventoryItemOfUser)
	}
}
