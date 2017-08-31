package product 

  import (
"github.com/hexya-erp/hexya/hexya/models/types"

) 

 func init() { 

 

pool.ProductPricelist().DeclareModel()
pool.ProductPricelist().Methods().GetDefaultCurrencyId().DeclareMethod(
`GetDefaultCurrencyId` ,
func (rs pool.ProductPricelistSet){
  /*def _get_default_currency_id(self):
        return self.env.user.company_id.currency_id.id

    */})
pool.ProductPricelist().Methods().GetDefaultItemIds().DeclareMethod(
`GetDefaultItemIds` ,
func (rs pool.ProductPricelistSet){
  /*def _get_default_item_ids(self):
        ProductPricelistItem = self.env['product.pricelist.item']
        vals = ProductPricelistItem.default_get(ProductPricelistItem._*/})
pool.ProductPricelist().AddCharField("Name", models.StringFieldParams{String :"Pricelist Name", Required : true, Translate: true})
pool.ProductPricelist().AddBooleanField("Active", models.SimpleFieldParams{String :"Active", Default: func(models.Environment, models.FieldMap) interface{} {return true} ,Help :"If unchecked, it will allow you to hide the pricelist without removing it."})
pool.ProductPricelist().AddOne2ManyField("Items", models.ReverseFieldParams{String :"Pricelist Items" ,RelationModel : pool.ProductPricelistItem() ,ReverseFK : "Pricelist" , JSON : "item_ids", NoCopy: false, Default : func(models.Environment, models.FieldMap) interface{}{
/*_get_default_item_ids(self):
        ProductPricelistItem = self.env['product.pricelist.item']
        vals = ProductPricelistItem.default_get(ProductPricelistItem._*/return 0}})
pool.ProductPricelist().AddMany2OneField("Currency",models.ForeignKeyFieldParams{String :"Currency" , RelationModel: pool.Currency() , JSON : "currency_id", Default : func(models.Environment, models.FieldMap) interface{}{
/*_get_default_currency_id(self):
        return self.env.user.company_id.currency_id.id

    */
return 0}, Required : true})
pool.ProductPricelist().AddMany2OneField("Company",models.ForeignKeyFieldParams{String :"Company" , RelationModel: pool.Company() , JSON : "company_id"})
pool.ProductPricelist().AddIntegerField("Sequence", models.SimpleFieldParams{String :"Sequence", Default: func(models.Environment, models.FieldMap) interface{} {return 16}})
pool.ProductPricelist().AddMany2ManyField("CountryGroups", models.Many2ManyFieldParams{String :"Country Groups" , RelationModel: pool.CountryGroup() , JSON : "country_group_ids"})
pool.ProductPricelist().Methods().NameGet().DeclareMethod(
`NameGet` ,
func (rs pool.ProductPricelistSet){
  //@api.multi
  /*def name_get(self):
        return [(pricelist.id, '%s (%s)' % (pricelist.name, pricelist.currency_id.name)) for pricelist in self]

    */})
pool.ProductPricelist().Methods().NameSearch().DeclareMethod(
`NameSearch` ,
func (rs pool.ProductPricelistSet , args struct{Args interface{}
Operator interface{}
Limit interface{}
}){
  //@api.model
  /*def name_search(self, name, args=None, operator='ilike', limit=100):
        if name and operator == '=' and not args:
            # search on the name of the pricelist and its currency, opposite of name_get(),
            # Used by the magic context filter in the product search view.
            query_args = {'name': name, 'limit': limit, 'lang': self._context.get('lang') or 'en_US'}
            query = """SELECT p.id
                       FROM ((
                                SELECT pr.id, pr.name
                                FROM product_pricelist pr JOIN
                                     res_currency cur ON 
                                         (pr.currency_id = cur.id)
                                WHERE pr.name || ' (' || cur.name || ')' = %(name)s
                            )
                            UNION (
                                SELECT tr.res_id as id, tr.value as name
                                FROM ir_translation tr JOIN
                                     product_pricelist pr ON (
                                        pr.id = tr.res_id AND
                                        tr.type = 'model' AND
                                        tr.name = 'product.pricelist,name' AND
                                        tr.lang = %(lang)s
                                     ) JOIN
                                     res_currency cur ON 
                                         (pr.currency_id = cur.id)
                                WHERE tr.value || ' (' || cur.name || ')' = %(name)s
                            )
                        ) p
                       ORDER BY p.name"""
            if limit:
                query += " LIMIT %(limit)s"
            self._cr.execute(query, query_args)
            ids = [r[0] for r in self._cr.fetchall()]
            # regular search() to apply ACLs - may limit results below limit in some cases
            pricelists = self.search([('id', 'in', ids)], limit=limit)
            if pricelists:
                return pricelists.name_get()
        return super(Pricelist, self).name_search(name, args, operator=operator, limit=limit)

    */})
pool.ProductPricelist().Methods().ComputePriceRuleMulti().DeclareMethod(
`ComputePriceRuleMulti` ,
func (rs pool.ProductPricelistSet , args struct{ProductsQtyPartner interface{}
Date interface{}
UomId interface{}
}){
  /*def _compute_price_rule_multi(self, products_qty_partner, date=False, uom_id=False):
        """ Low-level method - Multi pricelist, multi products
        Returns: dict{product_id: dict{pricelist_id: (price, suitable_rule)} }"""
        if not self.ids:
            pricelists = self.search([])
        else:
            pricelists = self
        results = {}
        for pricelist in pricelists:
            subres = pricelist._compute_price_rule(products_qty_partner, date=date, uom_id=uom_id)
            for product_id, price in subres.items():
                results.setdefault(product_id, {})
                results[product_id][pricelist.id] = price
        return results

    */})
pool.ProductPricelist().Methods().ComputePriceRule().DeclareMethod(
`ComputePriceRule` ,
func (rs pool.ProductPricelistSet , args struct{ProductsQtyPartner interface{}
Date interface{}
UomId interface{}
}){
  //@api.multi
  /*def _compute_price_rule(self, products_qty_partner, date=False, uom_id=False):
        """ Low-level method - Mono pricelist, multi products
        Returns: dict{product_id: (price, suitable_rule) for the given pricelist}

        If date in context: Date of the pricelist (%Y-%m-%d)

            :param products_qty_partner: list of typles products, quantity, partner
            :param datetime date: validity date
            :param ID uom_id: intermediate unit of measure
        """
        self.ensure_one()
        if not date:
            date = self._context.get('date', */})
pool.ProductPricelist().Methods().GetProductsPrice().DeclareMethod(
`GetProductsPrice` ,
func (rs pool.ProductPricelistSet , args struct{Products interface{}
Quantities interface{}
Partners interface{}
Date interface{}
UomId interface{}
}){
  /*def get_products_price(self, products, quantities, partners, date=False, uom_id=False):
        """ For a given pricelist, return price for products
        Returns: dict{product_id: product price}, in the given pricelist """
        self.ensure_one()
        return dict((product_id, res_tuple[0]) for product_id, res_tuple in self._compute_price_rule(zip(products, quantities, partners), date=date, uom_id=uom_id).iteritems())

    */})
pool.ProductPricelist().Methods().GetProductPrice().DeclareMethod(
`GetProductPrice` ,
func (rs pool.ProductPricelistSet , args struct{Product interface{}
Quantity interface{}
Partner interface{}
Date interface{}
UomId interface{}
}){
  /*def get_product_price(self, product, quantity, partner, date=False, uom_id=False):
        """ For a given pricelist, return price for a given product """
        self.ensure_one()
        return self._compute_price_rule([(product, quantity, partner)], date=date, uom_id=uom_id)[product.id][0]

    */})
pool.ProductPricelist().Methods().GetProductPriceRule().DeclareMethod(
`GetProductPriceRule` ,
func (rs pool.ProductPricelistSet , args struct{Product interface{}
Quantity interface{}
Partner interface{}
Date interface{}
UomId interface{}
}){
  /*def get_product_price_rule(self, product, quantity, partner, date=False, uom_id=False):
        """ For a given pricelist, return price and rule for a given product """
        self.ensure_one()
        return self._compute_price_rule([(product, quantity, partner)], date=date, uom_id=uom_id)[product.id]

    # Compatibility to remove after v10 - DEPRECATED
    */})
pool.ProductPricelist().Methods().PriceRuleGetMulti().DeclareMethod(
`PriceRuleGetMulti` ,
func (rs pool.ProductPricelistSet , args struct{Pricelist interface{}
ProductsByQtyByPartner interface{}
}){
  //@api.model
  /*def _price_rule_get_multi(self, pricelist, products_by_qty_by_partner):
        """ Low level method computing the result tuple for a given pricelist and multi products - return tuple """
        return pricelist._compute_price_rule(products_by_qty_by_partner)

    */})
pool.ProductPricelist().Methods().PriceGet().DeclareMethod(
`PriceGet` ,
func (rs pool.ProductPricelistSet , args struct{ProdId interface{}
Partner interface{}
}){
  //@api.multi
  /*def price_get(self, prod_id, qty, partner=None):
        """ Multi pricelist, mono product - returns price per pricelist """
        return dict((key, price[0]) for key, price in self.price_rule_get(prod_id, qty, partner=partner).items())

    */})
pool.ProductPricelist().Methods().PriceRuleGetMulti().DeclareMethod(
`PriceRuleGetMulti` ,
func (rs pool.ProductPricelistSet , args struct{ProductsByQtyByPartner interface{}
}){
  //@api.multi
  /*def price_rule_get_multi(self, products_by_qty_by_partner):
        """ Multi pricelist, multi product  - return tuple """
        return self._compute_price_rule_multi(products_by_qty_by_partner)

    */})
pool.ProductPricelist().Methods().PriceRuleGet().DeclareMethod(
`PriceRuleGet` ,
func (rs pool.ProductPricelistSet , args struct{ProdId interface{}
Partner interface{}
}){
  //@api.multi
  /*def price_rule_get(self, prod_id, qty, partner=None):
        """ Multi pricelist, mono product - return tuple """
        product = self.env['product.product'].browse([prod_id])
        return self._compute_price_rule_multi([(product, qty, partner)])[prod_id]

    */})
pool.ProductPricelist().Methods().PriceGetMulti().DeclareMethod(
`PriceGetMulti` ,
func (rs pool.ProductPricelistSet , args struct{Pricelist interface{}
ProductsByQtyByPartner interface{}
}){
  //@api.model
  /*def _price_get_multi(self, pricelist, products_by_qty_by_partner):
        """ Mono pricelist, multi product - return price per product """
        return pricelist.get_products_price(zip(**products_by_qty_by_partner))

    */})
pool.ProductPricelist().Methods().GetPartnerPricelist().DeclareMethod(
`GetPartnerPricelist` ,
func (rs pool.ProductPricelistSet , args struct{PartnerId interface{}
CompanyId interface{}
}){
  /*def _get_partner_pricelist(self, partner_id, company_id=None):
        """ Retrieve the applicable pricelist for a given partner in a given company.

            :param company_id: if passed, used for looking up properties,
             instead of current user's company
        """
        Partner = self.env['res.partner']
        Property = self.env['ir.property'].with_context(force_company=company_id or self.env.user.company_id.id)

        p = Partner.browse(partner_id)
        pl = Property.get('property_product_pricelist', Partner._name, '%s,%s' % (Partner._name, p.id))
        if pl:
            pl = pl[0].id

        if not pl:
            if p.country_id.code:
                pls = self.env['product.pricelist'].search([('country_group_ids.country_ids.code', '=', p.country_id.code)], limit=1)
                pl = pls and pls[0].id

        if not pl:
            # search pl where no country
            pls = self.env['product.pricelist'].search([('country_group_ids', '=', False)], limit=1)
            pl = pls and pls[0].id

        if not pl:
            prop = Property.get('property_product_pricelist', 'res.partner')
            pl = prop and prop[0].id

        if not pl:
            pls = self.env['product.pricelist'].search([], limit=1)
            pl = pls and pls[0].id

        return pl


*/})


pool.CountryGroup().DeclareModel()
pool.CountryGroup().AddMany2ManyField("Pricelists", models.Many2ManyFieldParams{String :"Pricelists" , RelationModel: pool.ProductPricelist() , JSON : "pricelist_ids"})


pool.ProductPricelistItem().DeclareModel()
pool.ProductPricelistItem().AddMany2OneField("ProductTmpl",models.ForeignKeyFieldParams{String :"Product Template" , RelationModel: pool.ProductTemplate() , JSON : "product_tmpl_id", OnDelete : models.Cascade , Help :"Specify a template if this rule only applies to one product template. Keep empty otherwise."})
pool.ProductPricelistItem().AddMany2OneField("Product",models.ForeignKeyFieldParams{String :"Product" , RelationModel: pool.ProductProduct() , JSON : "product_id", OnDelete : models.Cascade , Help :"Specify a product if this rule only applies to one product. Keep empty otherwise."})
pool.ProductPricelistItem().AddMany2OneField("Categ",models.ForeignKeyFieldParams{String :"Product Category" , RelationModel: pool.ProductCategory() , JSON : "categ_id", OnDelete : models.Cascade , Help :"Specify a product category if this rule only applies to products belonging to this category or its children categories. Keep empty otherwise."})
pool.ProductPricelistItem().AddIntegerField("MinQuantity", models.SimpleFieldParams{String :"Min. Quantity", Default: func(models.Environment, models.FieldMap) interface{} {return 1} ,Help :"For the rule to apply, bought/sold quantity must be greater  than or equal to the minimum quantity specified in this field.\n Expressed in the default unit of measure of the product."})
pool.ProductPricelistItem().AddSelectionField("AppliedOn", models.SelectionFieldParams{String :"Apply On", Selection : types.Selection{
"3_global" : "Global",
"2_product_category" : " Product Category",
"1_product" : "Product",
"0_product_variant" : "Product Variant",
}, Default: func(models.Environment, models.FieldMap) interface{} {return "3_global"}, Required : true, Help : "Pricelist Item applicable on selected option"})
pool.ProductPricelistItem().AddIntegerField("Sequence", models.SimpleFieldParams{String :"Sequence", Default: func(models.Environment, models.FieldMap) interface{} {return 5}, Required : true ,Help :"Gives the order in which the pricelist items will be checked. The evaluation gives highest priority to lowest sequence and stops as soon as a matching item is found."})
pool.ProductPricelistItem().AddSelectionField("Base", models.SelectionFieldParams{String :"Based on", Selection : types.Selection{
"list_price" : "Public Price",
"standard_price" : "Cost",
"pricelist" : "Other Pricelist",
}, Default: func(models.Environment, models.FieldMap) interface{} {return "list_price"}, Required : true, Help : "Base price for computation.\n' 'Public Price: The base price will be the Sale/public Price.\n' 'Cost Price : The base price will be the cost price.\n' 'Other Pricelist : Computation of the base price based on another Pricelist."})
pool.ProductPricelistItem().AddMany2OneField("BasePricelist",models.ForeignKeyFieldParams{String :"Other Pricelist" , RelationModel: pool.ProductPricelist() , JSON : "base_pricelist_id"})
pool.ProductPricelistItem().AddMany2OneField("Pricelist",models.ForeignKeyFieldParams{String :"Pricelist" , RelationModel: pool.ProductPricelist() , JSON : "pricelist_id", Index: true, OnDelete : models.Cascade})
pool.ProductPricelistItem().AddFloatField("PriceSurcharge", models.FloatFieldParams{String :"Price Surcharge"})
pool.ProductPricelistItem().AddFloatField("PriceDiscount", models.FloatFieldParams{String :"Price Discount", Default: func(models.Environment, models.FieldMap) interface{} {return 0}})
pool.ProductPricelistItem().AddFloatField("PriceRound", models.FloatFieldParams{String :"Price Rounding"})
pool.ProductPricelistItem().AddFloatField("PriceMinMargin", models.FloatFieldParams{String :"Min. Price Margin"})
pool.ProductPricelistItem().AddFloatField("PriceMaxMargin", models.FloatFieldParams{String :"Max. Price Margin"})
pool.ProductPricelistItem().AddMany2OneField("Company",models.ForeignKeyFieldParams{String :"Company" , RelationModel: pool.Company() , JSON : "company_id", Related: "Pricelist.Company"})
pool.ProductPricelistItem().Fields().Company().RevokeAccess(security.GroupEveryone, security.Write)
pool.ProductPricelistItem().AddMany2OneField("Currency",models.ForeignKeyFieldParams{String :"Currency" , RelationModel: pool.Currency() , JSON : "currency_id", Related: "Pricelist.Currency"})
pool.ProductPricelistItem().Fields().Currency().RevokeAccess(security.GroupEveryone, security.Write)
pool.ProductPricelistItem().AddDateField("DateStart", models.SimpleFieldParams{String :"Start Date" ,Help :"Starting date for the pricelist item validation"})
pool.ProductPricelistItem().AddDateField("DateEnd", models.SimpleFieldParams{String :"End Date" ,Help :"Ending valid for the pricelist item validation"})
pool.ProductPricelistItem().AddSelectionField("ComputePrice", models.SelectionFieldParams{String :"ComputePrice", Selection : types.Selection{
"fixed" : "Fix Price",
"percentage" : "Percentage (discount)",
"formula" : "Formula",
}, Index: true, Default: func(models.Environment, models.FieldMap) interface{} {return "fixed"}})
pool.ProductPricelistItem().AddFloatField("FixedPrice", models.FloatFieldParams{String :"Fixed Price"})
pool.ProductPricelistItem().AddFloatField("PercentPrice", models.FloatFieldParams{String :"Percentage Price')"})
pool.ProductPricelistItem().AddCharField("Name", models.StringFieldParams{String :"Name" ,Help :"Explicit rule name for this pricelist line."})
pool.ProductPricelistItem().AddCharField("Price", models.StringFieldParams{String :"Price" ,Help :"Explicit rule name for this pricelist line."})
pool.ProductPricelistItem().Methods().CheckRecursion().DeclareMethod(
`CheckRecursion` ,
func (rs pool.ProductPricelistItemSet){
  //@api.constrains('base_pricelist_id','pricelist_id','base')
  /*def _check_recursion(self):
        if any(item.base == 'pricelist' and item.pricelist_id and item.pricelist_id == item.base_pricelist_id for item in self):
            raise ValidationError(_('Error! You cannot assign the Main Pricelist as Other Pricelist in PriceList Item!'))
        return True

    */})
pool.ProductPricelistItem().Methods().CheckMargin().DeclareMethod(
`CheckMargin` ,
func (rs pool.ProductPricelistItemSet){
  //@api.constrains('price_min_margin','price_max_margin')
  /*def _check_margin(self):
        if any(item.price_min_margin > item.price_max_margin for item in self):
            raise ValidationError(_('Error! The minimum margin should be lower than the maximum margin.'))
        return True

    */})
pool.ProductPricelistItem().Methods().GetPricelistItemNamePrice().DeclareMethod(
`GetPricelistItemNamePrice` ,
func (rs pool.ProductPricelistItemSet){
  /*def _get_pricelist_item_name_price(self):
        if self.categ_id:
            self.name = _("Category: %s") % (self.categ_id.name)
        elif self.product_tmpl_id:
            self.name = self.product_tmpl_id.name
        elif self.product_id:
            self.name = self.product_id.display_name.replace('[%s]' % self.product_id.code, '')
        else:
            self.name = _("All Products")

        if self.compute_price == 'fixed':
            self.price = ("%s %s") % (self.fixed_price, self.pricelist_id.currency_id.name)
        elif self.compute_price == 'percentage':
            self.price = _("%s %% discount") % (self.percent_price)
        else:
            self.price = _("%s %% discount and %s surcharge") % (abs(self.price_discount), self.price_surcharge)

    */})
pool.ProductPricelistItem().Methods().OnchangeAppliedOn().DeclareMethod(
`OnchangeAppliedOn` ,
func (rs pool.ProductPricelistItemSet){
  //@api.onchange('applied_on')
  /*def _onchange_applied_on(self):
        if self.applied_on != '0_product_variant':
            self.product_id = False
        if self.applied_on != '1_product':
            self.product_tmpl_id = False
        if self.applied_on != '2_product_category':
            self.categ_id = False

    */})
pool.ProductPricelistItem().Methods().OnchangeComputePrice().DeclareMethod(
`OnchangeComputePrice` ,
func (rs pool.ProductPricelistItemSet){
  //@api.onchange('compute_price')
  /*def _onchange_compute_price(self):
        if self.compute_price != 'fixed':
            self.fixed_price = 0.0
        if self.compute_price != 'percentage':
            self.percent_price = 0.0
        if self.compute_price != 'formula':
            self.update({
                'price_discount': 0.0,
                'price_surcharge': 0.0,
                'price_round': 0.0,
                'price_min_margin': 0.0,
                'price_max_margin': 0.0,
            })

*/})
 
 }