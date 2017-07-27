package product

import (
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/models/types"

	"github.com/hexya-erp/hexya/pool"

	"github.com/hexya-erp/hexya/hexya/models"
)

func init() {

	pool.Pricelist().DeclareModel()
	pool.Pricelist().Method().GetDefaultCurrencyId().DeclareMethod(
		`GetDefaultCurrencyId`,
		func() { //def _get_default_currency_id(self):
		})
	pool.Pricelist().Method().GetDefaultItemIds().DeclareMethod(
		`GetDefaultItemIds`,
		func() { //def _get_default_item_ids(self):
			//ProductPricelistItem = self.env['product.pricelist.item']
			//vals = ProductPricelistItem.default_get(ProductPricelistItem._fields.keys())
			//vals.update(compute_price='formula')
			//return [[0, False, vals]]
			//name = fields.Char('Pricelist Name', required=True, translate=True)
			//active = fields.Boolean('Active', default=True, help="If unchecked, it will allow you to hide the pricelist without removing it.")
			//item_ids = fields.One2many(
			//'product.pricelist.item', 'pricelist_id', 'Pricelist Items',
			//copy=True, default=_get_default_item_ids)
			//currency_id = fields.Many2one('res.currency', 'Currency', default=_get_default_currency_id, required=True)
			//company_id = fields.Many2one('res.company', 'Company')
			//sequence = fields.Integer(default=16)
			//country_group_ids = fields.Many2many('res.country.group', 'res_country_group_pricelist_rel',
			//'pricelist_id', 'res_country_group_id', string='Country Groups')
		})
	pool.Pricelist().AddCharField("Name", models.StringFieldParams{String: "Name", Required: true, Translate: true})
	pool.Pricelist().AddBooleanField("Active", models.SimpleFieldParams{String: "Active", Default: func(models.Environment, models.FieldMap) interface{} { return true }, Help: "If unchecked, it will allow you to hide the pricelist without removing it."})
	pool.Pricelist().AddOne2ManyField("ItemIds", models.ReverseFieldParams{String: "product.pricelist.item", NoCopy: false, Default: func(models.Environment, models.FieldMap) interface{} { return _get_default_item_ids }})
	pool.Pricelist().AddMany2OneField("CurrencyId", models.ForeignKeyFieldParams{String: "Currency", RelationModel: pool.Res.Currency(), Default: func(models.Environment, models.FieldMap) interface{} { return _get_default_currency_id }, Required: true})
	pool.Pricelist().AddMany2OneField("CompanyId", models.ForeignKeyFieldParams{String: "Company')", RelationModel: pool.Res.Company()})
	pool.Pricelist().AddIntegerField("Sequence", models.SimpleFieldParams{String: "Sequence", Default: func(models.Environment, models.FieldMap) interface{} { return 16 }})
	pool.Pricelist().AddMany2ManyField("CountryGroupIds", models.Many2ManyFieldParams{String: "Country Groups", RelationModel: pool.Res.Country.Group()})
	pool.Pricelist().Method().NameGet().DeclareMethod(
		`NameGet`,
		func() { //def name_get(self):
			//return [(pricelist.id, '%s (%s)' % (pricelist.name, pricelist.currency_id.name)) for pricelist in self]
		})
	pool.Pricelist().Method().NameSearch().DeclareMethod(
		`NameSearch`,
		func() { //def name_search(self, name, args=None, operator='ilike', limit=100):
			//if name and operator == '=' and not args:
			//# search on the name of the pricelist and its currency, opposite of name_get(),
			//# Used by the magic context filter in the product search view.
			//query_args = {'name': name, 'limit': limit, 'lang': self._context.get('lang') or 'en_US'}
			//query = """SELECT p.id
			//FROM ((
			//SELECT pr.id, pr.name
			//FROM product_pricelist pr JOIN
			//res_currency cur ON
			//(pr.currency_id = cur.id)
			//WHERE pr.name || ' (' || cur.name || ')' = %(name)s
			//)
			//UNION (
			//SELECT tr.res_id as id, tr.value as name
			//FROM ir_translation tr JOIN
			//product_pricelist pr ON (
			//pr.id = tr.res_id AND
			//tr.type = 'model' AND
			//tr.name = 'product.pricelist,name' AND
			//tr.lang = %(lang)s
			//) JOIN
			//res_currency cur ON
			//(pr.currency_id = cur.id)
			//WHERE tr.value || ' (' || cur.name || ')' = %(name)s
			//)
			//) p
			//ORDER BY p.name"""
			//if limit:
			//query += " LIMIT %(limit)s"
			//self._cr.execute(query, query_args)
			//ids = [r[0] for r in self._cr.fetchall()]
			//# regular search() to apply ACLs - may limit results below limit in some cases
			//pricelists = self.search([('id', 'in', ids)], limit=limit)
			//if pricelists:
			//return pricelists.name_get()
		})
	pool.Pricelist().Method().ComputePriceRuleMulti().DeclareMethod(
		`ComputePriceRuleMulti`,
		func() { //def _compute_price_rule_multi(self, products_qty_partner, date=False, uom_id=False):
			//""" Low-level method - Multi pricelist, multi products
			//Returns: dict{product_id: dict{pricelist_id: (price, suitable_rule)} }"""
			//if not self.ids:
			//pricelists = self.search([])
			//else:
			//pricelists = self
			//results = {}
			//for pricelist in pricelists:
			//subres = pricelist._compute_price_rule(products_qty_partner, date=date, uom_id=uom_id)
			//for product_id, price in subres.items():
			//results.setdefault(product_id, {})
			//results[product_id][pricelist.id] = price
			//return results
		})
	pool.Pricelist().Method().ComputePriceRule().DeclareMethod(
		`ComputePriceRule`,
		func() { //def _compute_price_rule(self, products_qty_partner, date=False, uom_id=False):
			//""" Low-level method - Mono pricelist, multi products
			//Returns: dict{product_id: (price, suitable_rule) for the given pricelist}
			//If date in context: Date of the pricelist (%Y-%m-%d)
			//:param products_qty_partner: list of typles products, quantity, partner
			//:param datetime date: validity date
			//:param ID uom_id: intermediate unit of measure
			//"""
			//self.ensure_one()
			//if not date:
			//date = self._context.get('date', fields.Date.today())
			//if not uom_id and self._context.get('uom'):
			//uom_id = self._context['uom']
			//if uom_id:
			//# rebrowse with uom if given
			//product_ids = [item[0].id for item in products_qty_partner]
			//products = self.env['product.product'].with_context(uom=uom_id).browse(product_ids)
			//products_qty_partner = [(products[index], data_struct[1], data_struct[2]) for index, data_struct in enumerate(products_qty_partner)]
			//else:
			//products = [item[0] for item in products_qty_partner]
			//if not products:
			//return {}
			//categ_ids = {}
			//for p in products:
			//categ = p.categ_id
			//while categ:
			//categ_ids[categ.id] = True
			//categ = categ.parent_id
			//categ_ids = categ_ids.keys()
			//is_product_template = products[0]._name == "product.template"
			//if is_product_template:
			//prod_tmpl_ids = [tmpl.id for tmpl in products]
			//# all variants of all products
			//prod_ids = [p.id for p in
			//list(chain.from_iterable([t.product_variant_ids for t in products]))]
			//else:
			//prod_ids = [product.id for product in products]
			//prod_tmpl_ids = [product.product_tmpl_id.id for product in products]
			//# Load all rules
			//self._cr.execute(
			//'SELECT item.id '
			//'FROM product_pricelist_item AS item '
			//'LEFT JOIN product_category AS categ '
			//'ON item.categ_id = categ.id '
			//'WHERE (item.product_tmpl_id IS NULL OR item.product_tmpl_id = any(%s))'
			//'AND (item.product_id IS NULL OR item.product_id = any(%s))'
			//'AND (item.categ_id IS NULL OR item.categ_id = any(%s)) '
			//'AND (item.pricelist_id = %s) '
			//'AND (item.date_start IS NULL OR item.date_start<=%s) '
			//'AND (item.date_end IS NULL OR item.date_end>=%s)'
			//'ORDER BY item.applied_on, item.min_quantity desc, categ.parent_left desc',
			//(prod_tmpl_ids, prod_ids, categ_ids, self.id, date, date))
			//item_ids = [x[0] for x in self._cr.fetchall()]
			//items = self.env['product.pricelist.item'].browse(item_ids)
			//results = {}
			//for product, qty, partner in products_qty_partner:
			//results[product.id] = 0.0
			//suitable_rule = False
			//# Final unit price is computed according to `qty` in the `qty_uom_id` UoM.
			//# An intermediary unit price may be computed according to a different UoM, in
			//# which case the price_uom_id contains that UoM.
			//# The final price will be converted to match `qty_uom_id`.
			//qty_uom_id = self._context.get('uom') or product.uom_id.id
			//price_uom_id = product.uom_id.id
			//qty_in_product_uom = qty
			//if qty_uom_id != product.uom_id.id:
			//try:
			//qty_in_product_uom = self.env['product.uom'].browse([self._context['uom']])._compute_quantity(qty, product.uom_id)
			//except UserError:
			//# Ignored - incompatible UoM in context, use default product UoM
			//pass
			//# if Public user try to access standard price from website sale, need to call price_compute.
			//# TDE SURPRISE: product can actually be a template
			//price = product.price_compute('list_price')[product.id]
			//price_uom = self.env['product.uom'].browse([qty_uom_id])
			//for rule in items:
			//if rule.min_quantity and qty_in_product_uom < rule.min_quantity:
			//continue
			//if is_product_template:
			//if rule.product_tmpl_id and product.id != rule.product_tmpl_id.id:
			//continue
			//if rule.product_id and not (product.product_variant_count == 1 and product.product_variant_id.id == rule.product_id.id):
			//# product rule acceptable on template if has only one variant
			//continue
			//else:
			//if rule.product_tmpl_id and product.product_tmpl_id.id != rule.product_tmpl_id.id:
			//continue
			//if rule.product_id and product.id != rule.product_id.id:
			//continue
			//if rule.categ_id:
			//cat = product.categ_id
			//while cat:
			//if cat.id == rule.categ_id.id:
			//break
			//cat = cat.parent_id
			//if not cat:
			//continue
			//if rule.base == 'pricelist' and rule.base_pricelist_id:
			//price_tmp = rule.base_pricelist_id._compute_price_rule([(product, qty, partner)])[product.id][0]  # TDE: 0 = price, 1 = rule
			//price = rule.base_pricelist_id.currency_id.compute(price_tmp, self.currency_id, round=False)
			//else:
			//# if base option is public price take sale price else cost price of product
			//# price_compute returns the price in the context UoM, i.e. qty_uom_id
			//price = product.price_compute(rule.base)[product.id]
			//convert_to_price_uom = (lambda price: product.uom_id._compute_price(price, price_uom))
			//if price is not False:
			//if rule.compute_price == 'fixed':
			//price = convert_to_price_uom(rule.fixed_price)
			//elif rule.compute_price == 'percentage':
			//price = (price - (price * (rule.percent_price / 100))) or 0.0
			//else:
			//# complete formula
			//price_limit = price
			//price = (price - (price * (rule.price_discount / 100))) or 0.0
			//if rule.price_round:
			//price = tools.float_round(price, precision_rounding=rule.price_round)
			//if rule.price_surcharge:
			//price_surcharge = convert_to_price_uom(rule.price_surcharge)
			//price += price_surcharge
			//if rule.price_min_margin:
			//price_min_margin = convert_to_price_uom(rule.price_min_margin)
			//price = max(price, price_limit + price_min_margin)
			//if rule.price_max_margin:
			//price_max_margin = convert_to_price_uom(rule.price_max_margin)
			//price = min(price, price_limit + price_max_margin)
			//suitable_rule = rule
			//break
			//# Final price conversion into pricelist currency
			//if suitable_rule and suitable_rule.compute_price != 'fixed' and suitable_rule.base != 'pricelist':
			//price = product.currency_id.compute(price, self.currency_id, round=False)
			//results[product.id] = (price, suitable_rule and suitable_rule.id or False)
			//return results
		})
	pool.Pricelist().Method().GetProductsPrice().DeclareMethod(
		`GetProductsPrice`,
		func() { //def get_products_price(self, products, quantities, partners, date=False, uom_id=False):
			//""" For a given pricelist, return price for products
			//Returns: dict{product_id: product price}, in the given pricelist """
			//self.ensure_one()
		})
	pool.Pricelist().Method().GetProductPrice().DeclareMethod(
		`GetProductPrice`,
		func() { //def get_product_price(self, product, quantity, partner, date=False, uom_id=False):
			//""" For a given pricelist, return price for a given product """
			//self.ensure_one()
		})
	pool.Pricelist().Method().GetProductPriceRule().DeclareMethod(
		`GetProductPriceRule`,
		func() { //def get_product_price_rule(self, product, quantity, partner, date=False, uom_id=False):
			//""" For a given pricelist, return price and rule for a given product """
			//self.ensure_one()
			//return self._compute_price_rule([(product, quantity, partner)], date=date, uom_id=uom_id)[product.id]
			//# Compatibility to remove after v10 - DEPRECATED
		})
	pool.Pricelist().Method().PriceRuleGetMulti().DeclareMethod(
		`PriceRuleGetMulti`,
		func() { //def _price_rule_get_multi(self, pricelist, products_by_qty_by_partner):
			//""" Low level method computing the result tuple for a given pricelist and multi products - return tuple """
			//return pricelist._compute_price_rule(products_by_qty_by_partner)
		})
	pool.Pricelist().Method().PriceGet().DeclareMethod(
		`PriceGet`,
		func() { //def price_get(self, prod_id, qty, partner=None):
			//""" Multi pricelist, mono product - returns price per pricelist """
			//return dict((key, price[0]) for key, price in self.price_rule_get(prod_id, qty, partner=partner).items())
		})
	pool.Pricelist().Method().PriceRuleGetMulti().DeclareMethod(
		`PriceRuleGetMulti`,
		func() { //def price_rule_get_multi(self, products_by_qty_by_partner):
			//""" Multi pricelist, multi product  - return tuple """
			//return self._compute_price_rule_multi(products_by_qty_by_partner)
		})
	pool.Pricelist().Method().PriceRuleGet().DeclareMethod(
		`PriceRuleGet`,
		func() { //def price_rule_get(self, prod_id, qty, partner=None):
			//""" Multi pricelist, mono product - return tuple """
			//product = self.env['product.product'].browse([prod_id])
			//return self._compute_price_rule_multi([(product, qty, partner)])[prod_id]
		})
	pool.Pricelist().Method().PriceGetMulti().DeclareMethod(
		`PriceGetMulti`,
		func() { //def _price_get_multi(self, pricelist, products_by_qty_by_partner):
			//""" Mono pricelist, multi product - return price per product """
		})
	pool.Pricelist().Method().GetPartnerPricelist().DeclareMethod(
		`GetPartnerPricelist`,
		func() { //def _get_partner_pricelist(self, partner_id, company_id=None):
			//""" Retrieve the applicable pricelist for a given partner in a given company.
			//:param company_id: if passed, used for looking up properties,
			//instead of current user's company
			//"""
			//Partner = self.env['res.partner']
			//Property = self.env['ir.property'].with_context(force_company=company_id or self.env.user.company_id.id)
			//p = Partner.browse(partner_id)
			//pl = Property.get('property_product_pricelist', Partner._name, '%s,%s' % (Partner._name, p.id))
			//if pl:
			//pl = pl[0].id
			//if not pl:
			//if p.country_id.code:
			//pls = self.env['product.pricelist'].search([('country_group_ids.country_ids.code', '=', p.country_id.code)], limit=1)
			//pl = pls and pls[0].id
			//if not pl:
			//# search pl where no country
			//pls = self.env['product.pricelist'].search([('country_group_ids', '=', False)], limit=1)
			//pl = pls and pls[0].id
			//if not pl:
			//prop = Property.get('property_product_pricelist', 'res.partner')
			//pl = prop and prop[0].id
			//if not pl:
			//pls = self.env['product.pricelist'].search([], limit=1)
			//pl = pls and pls[0].id
		})

	pool.ResCountryGroup().DeclareModel()
	pool.ResCountryGroup().AddMany2ManyField("PricelistIds", models.Many2ManyFieldParams{String: "Pricelists", RelationModel: pool.Product.Pricelist()})

	pool.PricelistItem().DeclareModel()
	pool.PricelistItem().AddMany2OneField("ProductTmplId", models.ForeignKeyFieldParams{String: "Product Template", RelationModel: pool.Product.Template(), OnDelete: models.Cascade, Help: "Specify a template if this rule only applies to one product template. Keep empty otherwise."})
	pool.PricelistItem().AddMany2OneField("ProductId", models.ForeignKeyFieldParams{String: "Product", RelationModel: pool.Product.Product(), OnDelete: models.Cascade, Help: "Specify a product if this rule only applies to one product. Keep empty otherwise."})
	pool.PricelistItem().AddMany2OneField("CategId", models.ForeignKeyFieldParams{String: "Product Category", RelationModel: pool.Product.Category(), OnDelete: models.Cascade, Help: "Specify a product category if this rule only applies to products belonging to this category or its children categories. Keep empty otherwise."})
	pool.PricelistItem().AddIntegerField("MinQuantity", models.SimpleFieldParams{String: "Min. Quantity", Default: func(models.Environment, models.FieldMap) interface{} { return 1 }, Help: "For the rule to apply, bought/sold quantity must be greater  than or equal to the minimum quantity specified in this field.\n Expressed in the default unit of measure of the product."})
	pool.PricelistItem().AddSelectionField("AppliedOn", models.SelectionFieldParams{String: "AppliedOn", Selection: types.Selection{
		"3_global":           "3Global",
		"2_product_category": "2ProductCategory",
		"1_product":          "1Product",
		"0_product_variant":  "0ProductVariant",
	}, Default: func(models.Environment, models.FieldMap) interface{} { return "3_global" }, Required: true, Help: "Pricelist Item applicable on selected option"})
	pool.PricelistItem().AddIntegerField("Sequence", models.SimpleFieldParams{String: "Sequence", Default: func(models.Environment, models.FieldMap) interface{} { return 5 }, Required: true, Help: "Gives the order in which the pricelist items will be checked. The evaluation gives highest priority to lowest sequence and stops as soon as a matching item is found."})
	pool.PricelistItem().AddSelectionField("Base", models.SelectionFieldParams{String: "Base", Selection: types.Selection{
		"list_price":     "ListPrice",
		"standard_price": "StandardPrice",
		"pricelist":      "Pricelist",
	}, Default: func(models.Environment, models.FieldMap) interface{} { return "list_price" }, Required: true, Help: "Base price for computation.\n' 'Public Price: The base price will be the Sale/public Price.\n' 'Cost Price : The base price will be the cost price.\n' 'Other Pricelist : Computation of the base price based on another Pricelist."})
	pool.PricelistItem().AddMany2OneField("BasePricelistId", models.ForeignKeyFieldParams{String: "Other Pricelist')", RelationModel: pool.Product.Pricelist()})
	pool.PricelistItem().AddMany2OneField("PricelistId", models.ForeignKeyFieldParams{String: "Pricelist", RelationModel: pool.Product.Pricelist(), Index: true, OnDelete: models.Cascade})
	pool.PricelistItem().AddFloatField("PriceSurcharge", models.FloatFieldParams{String: "Price Surcharge"})
	pool.PricelistItem().AddFloatField("PriceDiscount", models.FloatFieldParams{String: "PriceDiscount", Default: func(models.Environment, models.FieldMap) interface{} { return 0 }})
	pool.PricelistItem().AddFloatField("PriceRound", models.FloatFieldParams{String: "Price Rounding"})
	pool.PricelistItem().AddFloatField("PriceMinMargin", models.FloatFieldParams{String: "Min. Price Margin"})
	pool.PricelistItem().AddFloatField("PriceMaxMargin", models.FloatFieldParams{String: "Max. Price Margin"})
	pool.PricelistItem().AddMany2OneField("CompanyId", models.ForeignKeyFieldParams{String: "Company", RelationModel: pool.Res.Company(), Related: "PricelistId.CompanyId", Stored: true})
	pool.PricelistItem().Fields().CompanyId().RevokeAccess(security.GroupEveryone, security.Write)
	pool.PricelistItem().AddMany2OneField("CurrencyId", models.ForeignKeyFieldParams{String: "Currency", RelationModel: pool.Res.Currency(), Related: "PricelistId.CurrencyId", Stored: true})
	pool.PricelistItem().Fields().CurrencyId().RevokeAccess(security.GroupEveryone, security.Write)
	pool.PricelistItem().AddDateField("DateStart", models.SimpleFieldParams{String: "DateStart", Help: "Starting date for the pricelist item validation"})
	pool.PricelistItem().AddDateField("DateEnd", models.SimpleFieldParams{String: "DateEnd", Help: "Ending valid for the pricelist item validation"})
	pool.PricelistItem().AddSelectionField("ComputePrice", models.SelectionFieldParams{String: "ComputePrice", Selection: types.Selection{
		"fixed":      "Fixed",
		"percentage": "Percentage",
		"":           "",
		"formula":    "Formula",
	}, Index: true, Default: func(models.Environment, models.FieldMap) interface{} { return "fixed" }})
	pool.PricelistItem().AddFloatField("FixedPrice", models.FloatFieldParams{String: "FixedPrice"})
	pool.PricelistItem().AddFloatField("PercentPrice", models.FloatFieldParams{String: "PercentPrice"})
	pool.PricelistItem().AddCharField("Name", models.StringFieldParams{String: "Name", Compute: "GetPricelistItemNamePrice", Help: "Explicit rule name for this pricelist line."})
	pool.PricelistItem().AddCharField("Price", models.StringFieldParams{String: "Price", Compute: "GetPricelistItemNamePrice", Help: "Explicit rule name for this pricelist line."})
	pool.PricelistItem().Method().CheckRecursion().DeclareMethod(
		`CheckRecursion`,
		func() { //def _check_recursion(self):
			//if any(item.base == 'pricelist' and item.pricelist_id and item.pricelist_id == item.base_pricelist_id for item in self):
			//raise ValidationError(_('Error! You cannot assign the Main Pricelist as Other Pricelist in PriceList Item!'))
			//return True
		})
	pool.PricelistItem().Method().CheckMargin().DeclareMethod(
		`CheckMargin`,
		func() { //def _check_margin(self):
			//if any(item.price_min_margin > item.price_max_margin for item in self):
			//raise ValidationError(_('Error! The minimum margin should be lower than the maximum margin.'))
			//return True
			//@api.one
			//@api.depends('categ_id', 'product_tmpl_id', 'product_id', 'compute_price', 'fixed_price', \
		})
	pool.PricelistItem().Method().GetPricelistItemNamePrice().DeclareMethod(
		`GetPricelistItemNamePrice`,
		func() { //def _get_pricelist_item_name_price(self):
			//if self.categ_id:
			//self.name = _("Category: %s") % (self.categ_id.name)
			//elif self.product_tmpl_id:
			//self.name = self.product_tmpl_id.name
			//elif self.product_id:
			//self.name = self.product_id.display_name.replace('[%s]' % self.product_id.code, '')
			//else:
			//self.name = _("All Products")
			//if self.compute_price == 'fixed':
			//self.price = ("%s %s") % (self.fixed_price, self.pricelist_id.currency_id.name)
			//elif self.compute_price == 'percentage':
			//self.price = _("%s %% discount") % (self.percent_price)
			//else:
			//self.price = _("%s %% discount and %s surcharge") % (abs(self.price_discount), self.price_surcharge)
		})
	pool.PricelistItem().Method().OnchangeAppliedOn().DeclareMethod(
		`OnchangeAppliedOn`,
		func() { //def _onchange_applied_on(self):
			//if self.applied_on != '0_product_variant':
			//self.product_id = False
			//if self.applied_on != '1_product':
			//self.product_tmpl_id = False
			//if self.applied_on != '2_product_category':
			//self.categ_id = False
		})
	pool.PricelistItem().Method().OnchangeComputePrice().DeclareMethod(
		`OnchangeComputePrice`,
		func() { //def _onchange_compute_price(self):
			//if self.compute_price != 'fixed':
			//self.fixed_price = 0.0
			//if self.compute_price != 'percentage':
			//self.percent_price = 0.0
			//if self.compute_price != 'formula':
			//self.update({
			//'price_discount': 0.0,
			//'price_surcharge': 0.0,
			//'price_round': 0.0,
			//'price_min_margin': 0.0,
			//'price_max_margin': 0.0,
		})

}
