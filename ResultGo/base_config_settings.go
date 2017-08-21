package  

  import (
"github.com/hexya-erp/hexya/hexya/models/types"

) 

 func init() { 

 

pool.BaseConfigSettings().DeclareModel()
pool.BaseConfigSettings().AddBooleanField("CompanyShareProduct", models.SimpleFieldParams{String :"Share product to all companies" ,Help :"Share your product to all companies defined in your instance.\n  * Checked : Product are visible for every company, even if a company is defined on the partner.\n  * Unchecked : Each company can see only its product "})
pool.BaseConfigSettings().AddSelectionField("GroupProductVariant", models.SelectionFieldParams{String :"Product Variants", Selection : types.Selection{
"0" : "No variants on products",
"1" : "Products can have several attributes defining variants (Example: size color...)",
}, Help : "Work with product variant allows you to define some variant of the same products"})
pool.BaseConfigSettings().Methods().GetDefaultCompanyShareProduct().DeclareMethod(
`GetDefaultCompanyShareProduct` ,
func (rs pool.BaseConfigSettingsSet , args struct{Fields interface{}
}){
  //@api.model
  /*def get_default_company_share_product(self, fields):
        product_rule = self.env.ref('product.product_comp_rule')
        return {
            'company_share_product': not bool(product_rule.active)
        }

    */})
pool.BaseConfigSettings().Methods().SetAuthCompanyShareProduct().DeclareMethod(
`SetAuthCompanyShareProduct` ,
func (rs pool.BaseConfigSettingsSet){
  //@api.multi
  /*def set_auth_company_share_product(self):
        self.ensure_one()
        product_rule = self.env.ref('product.product_comp_rule')
        product_rule.write({'active': not bool(self.company_share_product)})
*/})
 
 }