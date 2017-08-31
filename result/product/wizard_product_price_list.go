package product 

  import (

) 

 func init() { 

 

pool.ProductPriceList().DeclareTransientModel()
pool.ProductPriceList().AddMany2OneField("PriceList",models.ForeignKeyFieldParams{String :"PriceList" , RelationModel: pool.ProductPricelist(), Required : true})
pool.ProductPriceList().AddIntegerField("Qty1", models.SimpleFieldParams{String :"Quantity-1", Default: func(models.Environment, models.FieldMap) interface{} {return 1}})
pool.ProductPriceList().AddIntegerField("Qty2", models.SimpleFieldParams{String :"Quantity-2", Default: func(models.Environment, models.FieldMap) interface{} {return 5}})
pool.ProductPriceList().AddIntegerField("Qty3", models.SimpleFieldParams{String :"Quantity-3", Default: func(models.Environment, models.FieldMap) interface{} {return 10}})
pool.ProductPriceList().AddIntegerField("Qty4", models.SimpleFieldParams{String :"Quantity-4", Default: func(models.Environment, models.FieldMap) interface{} {return 0}})
pool.ProductPriceList().AddIntegerField("Qty5", models.SimpleFieldParams{String :"Quantity-5", Default: func(models.Environment, models.FieldMap) interface{} {return 0}})
pool.ProductPriceList().Methods().PrintReport().DeclareMethod(
`PrintReport` ,
func (rs pool.ProductPriceListSet){
  //@api.multi
  /*def print_report(self):
        """
        To get the date and print the report
        @return : return report
        """
        datas = {'ids': self.env.context.get('active_ids', [])}
        res = self.read(['price_list', 'qty1', 'qty2', 'qty3', 'qty4', 'qty5'])
        res = res and res[0] or {}
        res['price_list'] = res['price_list'][0]
        datas['form'] = res
        return self.env['report'].get_action([], 'product.report_pricelist', data=datas)
*/})
 
 }