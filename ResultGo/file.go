package product 

  import (
"github.com/hexya-erp/hexya/hexya/models/types"

) 

 func init() { 

 

pool.ProductCategory().DeclareModel()
pool.ProductCategory().AddCharField("Name", models.StringFieldParams{String :"Name", Index: true, Required: true, Translate: true})
pool.ProductCategory().AddMany2OneField("ParentId",models.ForeignKeyFieldParams{String :"Parent Category" , RelationModel: pool.ProductCategory(), Index: true, OnDelete : models.Cascade})
pool.ProductCategory().AddOne2ManyField("ChildId", models.ReverseFieldParams{String :"Child Categories')" ,RelationModel : pool.ProductCategory ,ReverseFK : "ParentId"})
pool.ProductCategory().AddSelectionField("Type", models.SelectionFieldParams{String :"Category Type", Selection : types.Selection{
 "view" : "View",
 "normal" : "Normal",
}, Default: func(models.Environment, models.FieldMap) interface{} {return "normal"}, Help : "A category of the view type is a virtual category that can be used as the parent of another category to create a hierarchical structure."})
pool.ProductCategory().AddIntegerField("ParentLeft", models.SimpleFieldParams{String :"Left Parent", Index: true})
pool.ProductCategory().AddIntegerField("ParentRight", models.SimpleFieldParams{String :"Right Parent", Index: true})
pool.ProductCategory().AddIntegerField("ProductCount", models.SimpleFieldParams{String :"# Products", Compute : pool.ComputeProductCount() ,Help :"The number of products under this category "})


pool.ProductPriceHistory().DeclareModel()
pool.ProductPriceHistory().AddMany2OneField("CompanyId",models.ForeignKeyFieldParams{String :"Company" , RelationModel: pool.ResCompany(), Default : pool.ProductPriceHistory.GetDefaultCompanyId(), Required: true})
pool.ProductPriceHistory().AddMany2OneField("ProductId",models.ForeignKeyFieldParams{String :"Product" , RelationModel: pool.ProductProduct(), OnDelete : models.Cascade, Required: true})
pool.ProductPriceHistory().AddDateTimeField("Datetime", models.SimpleFieldParams{})
pool.ProductPriceHistory().AddFloatField("Cost", models.FloatFieldParams{String :"Cost"})


pool.ProductProduct().DeclareModel()
pool.ProductProduct().AddFloatField("Price", models.FloatFieldParams{String :"Price", Compute: "ComputeProductPrice"})
pool.ProductProduct().AddFloatField("PriceExtra", models.FloatFieldParams{String :"Variant Price Extra", Compute: "ComputeProductPriceExtra"})
pool.ProductProduct().AddFloatField("LstPrice", models.FloatFieldParams{String :"Sale Price", Compute: "ComputeProductLstPrice"})
pool.ProductProduct().AddCharField("DefaultCode", models.StringFieldParams{String :"Internal Reference", Index: true})
pool.ProductProduct().AddCharField("Code", models.StringFieldParams{String :"Internal Reference", Compute: "ComputeProductCode"})
pool.ProductProduct().AddCharField("PartnerRef", models.StringFieldParams{String :"Customer Ref", Compute: "ComputePartnerRef"})
pool.ProductProduct().AddBooleanField("Active", models.SimpleFieldParams{String :"Active", Default: func(models.Environment, models.FieldMap) interface{} {return true} ,Help :"If unchecked, it will allow you to hide the product without removing it."})
pool.ProductProduct().AddMany2OneField("ProductTmplId",models.ForeignKeyFieldParams{String :"Product Template" , RelationModel: pool.ProductTemplate(), Index: true, OnDelete : models.Cascade, Required: true})
pool.ProductProduct().AddCharField("Barcode", models.StringFieldParams{String :"Barcode", NoCopy: true ,Help :"International Article Number used for product identification." ,Help :"International Article Number used for product identification."})
pool.ProductProduct().AddMany2ManyField("AttributeValueIds", models.Many2ManyFieldParams{String :"Attributes" , RelationModel: pool.ProductAttributeValue(), OnDelete : models.Restrict})
pool.ProductProduct().AddBinaryField("ImageVariant", models.SimpleFieldParams{String :"Variant Image" ,Help :"This field holds the image used as image for the product variant, limited to 1024x1024px."})
pool.ProductProduct().AddBinaryField("ImageSmall", models.SimpleFieldParams{String :"Small-sized image", Compute: "ComputeImages" ,Help :"Image of the product variant "})
pool.ProductProduct().AddBinaryField("Image", models.SimpleFieldParams{String :"Big-sized image", Compute: "ComputeImages" ,Help :"Image of the product variant "})
pool.ProductProduct().AddBinaryField("ImageMedium", models.SimpleFieldParams{String :"Medium-sized image", Compute: "ComputeImages" ,Help :"Image of the product variant "})
pool.ProductProduct().AddFloatField("StandardPrice", models.FloatFieldParams{String :"Cost"})
pool.ProductProduct().AddFloatField("Volume", models.FloatFieldParams{String :"Volume" ,Help :"The volume in m3."})
pool.ProductProduct().AddFloatField("Weight", models.FloatFieldParams{String :"Weight"})
pool.ProductProduct().AddMany2ManyField("PricelistItemIds", models.Many2ManyFieldParams{String :"Pricelist Items" , RelationModel: pool.ProductPricelistItem(), Compute: "GetPricelistItems"})
pool.ProductProduct().AddSQLConstraint("BarcodeUniq" , "Unique(barcode)" , "A barcode can only be assigned to one product !")


pool.ProductPackaging().DeclareModel()
pool.ProductPackaging().AddCharField("Name", models.StringFieldParams{String :"Packaging Type", Required: true})
pool.ProductPackaging().AddIntegerField("Sequence", models.SimpleFieldParams{String :"Sequence", Default: func(models.Environment, models.FieldMap) interface{} {return 1} ,Help :"The first in the sequence is the default one."})
pool.ProductPackaging().AddMany2OneField("ProductTmplId",models.ForeignKeyFieldParams{String :"Product" , RelationModel: pool.ProductTemplate()})
pool.ProductPackaging().AddFloatField("Qty", models.FloatFieldParams{String :"Quantity per Package" ,Help :"The total number of products you can have per pallet or box."})


pool.SuppliferInfo().DeclareModel()
pool.SuppliferInfo().AddMany2OneField("Name",models.ForeignKeyFieldParams{String :"Vendor" , RelationModel: pool.ResPartner()})
pool.SuppliferInfo().AddCharField("ProductName", models.StringFieldParams{String :"Vendor Product Name" ,Help :"This vendor's product name will be used when printing a request for quotation. Keep empty to use the internal one." ,Help :"This vendor's product name will be used when printing a request for quotation. Keep empty to use the internal one."})
pool.SuppliferInfo().AddCharField("ProductCode", models.StringFieldParams{String :"Vendor Product Code" ,Help :"This vendor's product code will be used when printing a request for quotation. Keep empty to use the internal one." ,Help :"This vendor's product code will be used when printing a request for quotation. Keep empty to use the internal one."})
pool.SuppliferInfo().AddIntegerField("Sequence", models.SimpleFieldParams{String :"Sequence", Default: func(models.Environment, models.FieldMap) interface{} {return 1} ,Help :"Assigns the priority to the list of product vendor."})
pool.SuppliferInfo().AddMany2OneField("ProductUom",models.ForeignKeyFieldParams{String :"Vendor Unit of Measure" , RelationModel: pool.ProductUom(), Related: "ProductTmplIdUomPoId" , Help :"This comes from the product form."})
pool.SuppliferInfo().Fields().ProductUom().RevokeAccess(security.GroupEveryone, security.Write)
pool.SuppliferInfo().AddFloatField("MinQty", models.FloatFieldParams{String :"Minimal Quantity", Default: func(models.Environment, models.FieldMap) interface{} {return 0.0}, Required: true ,Help :"The minimal quantity to purchase from this vendor, expressed in the vendor Product Unit of Measure if not any, in the default unit of measure of the product otherwise."})
pool.SuppliferInfo().AddFloatField("Price", models.FloatFieldParams{String :"Price", Default: func(models.Environment, models.FieldMap) interface{} {return 0.0}})
pool.SuppliferInfo().AddMany2OneField("CompanyId",models.ForeignKeyFieldParams{String :"Company" , RelationModel: pool.ResCompany(), Default: func(models.Environment, models.FieldMap) interface{} {return lambda self: self.env.user.company_id.id}, Index: true})
pool.SuppliferInfo().AddMany2OneField("CurrencyId",models.ForeignKeyFieldParams{String :"Currency" , RelationModel: pool.ResCurrency(), Default: func(models.Environment, models.FieldMap) interface{} {return lambda self: self.env.user.company_id.currency_id.id}, Required: true})
pool.SuppliferInfo().AddDateField("DateStart", models.SimpleFieldParams{String :"Start Date" ,Help :"Start date for this vendor price"})
pool.SuppliferInfo().AddDateField("DateEnd", models.SimpleFieldParams{String :"End Date" ,Help :"End date for this vendor price"})
pool.SuppliferInfo().AddMany2OneField("ProductId",models.ForeignKeyFieldParams{String :"Product Variant" , RelationModel: pool.ProductProduct() , Help :"When this field is filled in, the vendor data will only apply to the variant."})
pool.SuppliferInfo().AddMany2OneField("ProductTmplId",models.ForeignKeyFieldParams{String :"Product Template" , RelationModel: pool.ProductTemplate(), Index: true, OnDelete : models.Cascade})
pool.SuppliferInfo().AddIntegerField("Delay", models.SimpleFieldParams{String :"Delivery Lead Time", Default: func(models.Environment, models.FieldMap) interface{} {return 1}, Required: true ,Help :"Lead time in days between the confirmation of the purchase order and the receipt of the products in your warehouse. Used by the scheduler for automatic computation of the purchase order planning."})
 
 }