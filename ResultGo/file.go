package product

import (
 "github.com/hexya-erp/hexya/pool"
 "github.com/hexya-erp/hexya/hexya/models"
)

func init() {

 

pool.BaseConfigSettings(models.T().DeclareModel()
pool.BaseConfigSettings(models.T().AddBooleanField("CompanyShareProduct", models.SimpleFieldParams{String :"Share product to all companies" ,Help :"Share your product to all companies defined in your instance.\n  * Checked : Product are visible for every company, even if a company is defined on the partner.\n  * Unchecked : Each company can see only its product "})
pool.BaseConfigSettings(models.T().AddSelectionField("GroupProductVariant", models.SelectionFieldParams{})
pool.BaseConfigSettings(models.T().Method().GetDefaultCompanyShareProduct().DeclareMethod(
`GetDefaultCompanyShareProduct` ,
func (){//def get_default_company_share_product(self, fields): 
//product_rule = self.env.ref('product.product_comp_rule') 
//return { 
//'company_share_product': not bool(product_rule.active) 
//} 
})
pool.BaseConfigSettings(models.T().Method().SetAuthCompanyShareProduct().DeclareMethod(
`SetAuthCompanyShareProduct` ,
func (){//def set_auth_company_share_product(self): 
//self.ensure_one() 
//product_rule = self.env.ref('product.product_comp_rule') 
})
 
 }