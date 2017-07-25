package product 

 func init() { 

 

pool.ProductCategory().DeclareModel()
pool.ProductCategory().AddCharField("Name", models.StringFieldParams{String :"Name", Index: true, Required: true, Translate: true})
pool.ProductCategory().AddMany2OneField("ParentId",models.ForeignKeyFieldParams{String :"Parent Category" , RelationModel: pool.Product.Category(), Index: true, OnDelete : models.Cascade})
pool.ProductCategory().AddOne2ManyField("ChildId", models.ReverseFieldParams{String :"Child Categories')"})
pool.ProductCategory().AddSelectionField("Type", models.SelectionFieldParams{String :"["})
pool.ProductCategory().AddIntegerField("ParentLeft", models.SimpleFieldParams{String :"Left Parent", Index: true})
pool.ProductCategory().AddIntegerField("ParentRight", models.SimpleFieldParams{String :"Right Parent", Index: true})
pool.ProductCategory().AddIntegerField("ProductCount", models.SimpleFieldParams{String :"# Products", Compute : pool.ComputeProductCount() ,Help :"The number of products under this category "})
pool.ProductCategory().Method().ComputeProductCount().DeclareMethod(
`ComputeProductCount` ,
func (){//def _compute_product_count(self): 
//read_group_res = self.env['product.template'].read_group([('categ_id', 'in', self.ids)], ['categ_id'], ['categ_id']) 
//group_data = dict((data['categ_id'][0], data['categ_id_count']) for data in read_group_res) 
//for categ in self: 
//categ.product_count = group_data.get(categ.id, 0) 
})
pool.ProductCategory().Method().CheckCategoryRecursion().DeclareMethod(
`CheckCategoryRecursion` ,
func (){//def _check_category_recursion(self): 
//if not self._check_recursion(): 
//raise ValidationError(_('Error ! You cannot create recursive categories.')) 
//return True 
})
pool.ProductCategory().Method().NameGet().DeclareMethod(
`NameGet` ,
func (){//def name_get(self): 
//def get_names(cat): 
//""" Return the list [cat.name, cat.parent_id.name, ...] """ 
//res = [] 
//while cat: 
//res.append(cat.name) 
//cat = cat.parent_id 
//return res 
//return [(cat.id, " / ".join(reversed(get_names(cat)))) for cat in self] 
})
pool.ProductCategory().Method().GetNames().DeclareMethod(
`GetNames` ,
func (){//def get_names(cat): 
//""" Return the list [cat.name, cat.parent_id.name, ...] """ 
//res = [] 
//while cat: 
//res.append(cat.name) 
//cat = cat.parent_id 
//return res 
//return [(cat.id, " / ".join(reversed(get_names(cat)))) for cat in self] 
})
pool.ProductCategory().Method().NameSearch().DeclareMethod(
`NameSearch` ,
func (){//def name_search(self, name, args=None, operator='ilike', limit=100): 
//if not args: 
//args = [] 
//if name: 
//# Be sure name_search is symetric to name_get 
//category_names = name.split(' / ') 
//parents = list(category_names) 
//child = parents.pop() 
//domain = [('name', operator, child)] 
//if parents: 
//names_ids = self.name_search(' / '.join(parents), args=args, operator='ilike', limit=limit) 
//category_ids = [name_id[0] for name_id in names_ids] 
//if operator in expression.NEGATIVE_TERM_OPERATORS: 
//categories = self.search([('id', 'not in', category_ids)]) 
//domain = expression.OR([[('parent_id', 'in', categories.ids)], domain]) 
//else: 
//domain = expression.AND([[('parent_id', 'in', category_ids)], domain]) 
//for i in range(1, len(category_names)): 
//domain = [[('name', operator, ' / '.join(category_names[-1 - i:]))], domain] 
//if operator in expression.NEGATIVE_TERM_OPERATORS: 
//domain = expression.AND(domain) 
//else: 
//domain = expression.OR(domain) 
//categories = self.search(expression.AND([domain, args]), limit=limit) 
//else: 
//categories = self.search(args, limit=limit) 
})


pool.ProductPriceHistory().DeclareModel()
pool.ProductPriceHistory().Method().GetDefaultCompanyId().DeclareMethod(
`GetDefaultCompanyId` ,
func (){//def _get_default_company_id(self): 
//return self._context.get('force_company', self.env.user.company_id.id) 
//company_id = fields.Many2one('res.company', string='Company', 
//default=_get_default_company_id, required=True) 
//product_id = fields.Many2one('product.product', 'Product', ondelete='cascade', required=True) 
//datetime = fields.Datetime('Date', default=fields.Datetime.now) 
})
pool.ProductPriceHistory().AddMany2OneField("CompanyId",models.ForeignKeyFieldParams{String :"Company" , RelationModel: pool.Res.Company(), Default: func(models.Environment, models.FieldMap) interface{} {return _get_default_company_id}, Required: true})
pool.ProductPriceHistory().AddMany2OneField("ProductId",models.ForeignKeyFieldParams{String :"Product" , RelationModel: pool.Product.Product(), OnDelete : models.Cascade, Required: true})
pool.ProductPriceHistory().AddDateTimeField("Datetime", models.SimpleFieldParams{})
pool.ProductPriceHistory().AddFloatField("Cost", models.FloatFieldParams{String :"Cost"})


pool.ProductProduct().DeclareModel()
pool.ProductProduct().AddFloatField("Price", models.FloatFieldParams{String :"Price", Compute: "ComputeProductPrice"})
pool.ProductProduct().AddFloatField("PriceExtra", models.FloatFieldParams{String :"Variant Price Extra", Compute: "ComputeProductPriceExtra"})
pool.ProductProduct().AddFloatField("LstPrice", models.FloatFieldParams{String :"Sale Price", Compute: "ComputeProductLstPrice"})
pool.ProductProduct().AddCharField("DefaultCode", models.StringFieldParams{String :"Internal Reference", Index: true})
pool.ProductProduct().AddCharField("Code", models.StringFieldParams{String :"Internal Reference", Compute: "ComputeProductCode"})
pool.ProductProduct().AddCharField("PartnerRef", models.StringFieldParams{String :"Customer Ref", Compute: "ComputePartnerRef"})
pool.ProductProduct().AddBooleanField("Active", models.SimpleFieldParams{String :"Active", Default: func(models.Environment, models.FieldMap) interface{} {return true} ,Help :"If unchecked, it will allow you to hide the product without removing it.)"})
pool.ProductProduct().AddMany2OneField("ProductTmplId",models.ForeignKeyFieldParams{String :"Product Template" , RelationModel: pool.Product.Template(), Index: true, OnDelete : models.Cascade, Required: true})
pool.ProductProduct().AddCharField("Barcode", models.StringFieldParams{String :"Barcode", NoCopy: true ,Help :"International Article Number used for product identification.)"})
pool.ProductProduct().AddMany2ManyField("AttributeValueIds", models.Many2ManyFieldParams{String :"Attributes" , RelationModel: pool.Product.Attribute.Value(), OnDelete : models.Restrict})
pool.ProductProduct().AddBinaryField("ImageVariant", models.SimpleFieldParams{String :"Variant Image" ,Help :"This field holds the image used as image for the product variant, limited to 1024x1024px.)"})
pool.ProductProduct().AddBinaryField("ImageSmall", models.SimpleFieldParams{String :"Small-sized image", Compute: "ComputeImages" ,Help :"Image of the product variant "})
pool.ProductProduct().AddBinaryField("Image", models.SimpleFieldParams{String :"Big-sized image", Compute: "ComputeImages" ,Help :"Image of the product variant "})
pool.ProductProduct().AddBinaryField("ImageMedium", models.SimpleFieldParams{String :"Medium-sized image", Compute: "ComputeImages" ,Help :"Image of the product variant "})
pool.ProductProduct().AddFloatField("StandardPrice", models.FloatFieldParams{String :"Cost"})
pool.ProductProduct().AddFloatField("Volume", models.FloatFieldParams{String :"Volume" ,Help :"The volume in m3.)"})
pool.ProductProduct().AddFloatField("Weight", models.FloatFieldParams{String :"Weight"})
pool.ProductProduct().AddMany2ManyField("PricelistItemIds", models.Many2ManyFieldParams{String :"Pricelist Items" , RelationModel: pool.Product.Pricelist.Item(), Compute: "GetPricelistItems"})
pool.ProductProduct().AddSQLConstraint("BarcodeUniq" , "Unique(barcode)" , ("A barcode can only be assigned to one product !"))
pool.ProductProduct().Method().ComputeProductPrice().DeclareMethod(
`ComputeProductPrice` ,
func (){//def _compute_product_price(self): 
//prices = {} 
//pricelist_id_or_name = self._context.get('pricelist') 
//if pricelist_id_or_name: 
//pricelist = None 
//partner = self._context.get('partner', False) 
//quantity = self._context.get('quantity', 1.0) 
//# Support context pricelists specified as display_name or ID for compatibility 
//if isinstance(pricelist_id_or_name, basestring): 
//pricelist_name_search = self.env['product.pricelist'].name_search(pricelist_id_or_name, operator='=', limit=1) 
//if pricelist_name_search: 
//pricelist = self.env['product.pricelist'].browse([pricelist_name_search[0][0]]) 
//elif isinstance(pricelist_id_or_name, (int, long)): 
//pricelist = self.env['product.pricelist'].browse(pricelist_id_or_name) 
//if pricelist: 
//quantities = [quantity] * len(self) 
//partners = [partner] * len(self) 
//prices = pricelist.get_products_price(self, quantities, partners) 
//for product in self: 
})
pool.ProductProduct().Method().SetProductPrice().DeclareMethod(
`SetProductPrice` ,
func (){//def _set_product_price(self): 
//for product in self: 
//if self._context.get('uom'): 
//value = self.env['product.uom'].browse(self._context['uom'])._compute_price(product.price, product.uom_id) 
//else: 
//value = product.price 
//value -= product.price_extra 
})
pool.ProductProduct().Method().SetProductLstPrice().DeclareMethod(
`SetProductLstPrice` ,
func (){//def _set_product_lst_price(self): 
//for product in self: 
//if self._context.get('uom'): 
//value = self.env['product.uom'].browse(self._context['uom'])._compute_price(product.lst_price, product.uom_id) 
//else: 
//value = product.lst_price 
//value -= product.price_extra 
//product.write({'list_price': value}) 
})
pool.ProductProduct().Method().ComputeProductPriceExtra().DeclareMethod(
`ComputeProductPriceExtra` ,
func (){//def _compute_product_price_extra(self): 
//# TDE FIXME: do a real multi and optimize a bit ? 
//for product in self: 
//price_extra = 0.0 
//for attribute_price in product.mapped('attribute_value_ids.price_ids'): 
//if attribute_price.product_tmpl_id == product.product_tmpl_id: 
//price_extra += attribute_price.price_extra 
//product.price_extra = price_extra 
})
pool.ProductProduct().Method().ComputeProductLstPrice().DeclareMethod(
`ComputeProductLstPrice` ,
func (){//def _compute_product_lst_price(self): 
//to_uom = None 
//if 'uom' in self._context: 
//to_uom = self.env['product.uom'].browse([self._context['uom']]) 
//for product in self: 
//if to_uom: 
//list_price = product.uom_id._compute_price(product.list_price, to_uom) 
//else: 
//list_price = product.list_price 
//product.lst_price = list_price + product.price_extra 
})
pool.ProductProduct().Method().ComputeProductCode().DeclareMethod(
`ComputeProductCode` ,
func (){//def _compute_product_code(self): 
//for supplier_info in self.seller_ids: 
//if supplier_info.name.id == self._context.get('partner_id'): 
//self.code = supplier_info.product_code or self.default_code 
//else: 
//self.code = self.default_code 
})
pool.ProductProduct().Method().ComputePartnerRef().DeclareMethod(
`ComputePartnerRef` ,
func (){//def _compute_partner_ref(self): 
//for supplier_info in self.seller_ids: 
//if supplier_info.name.id == self._context.get('partner_id'): 
//product_name = supplier_info.product_name or self.default_code 
//else: 
//product_name = self.name 
//self.partner_ref = '%s%s' % (self.code and '[%s] ' % self.code or '', product_name) 
//@api.one 
})
pool.ProductProduct().Method().ComputeImages().DeclareMethod(
`ComputeImages` ,
func (){//def _compute_images(self): 
//if self._context.get('bin_size'): 
//self.image_medium = self.image_variant 
//self.image_small = self.image_variant 
//self.image = self.image_variant 
//else: 
//resized_images = tools.image_get_resized_images(self.image_variant, return_big=True, avoid_resize_medium=True) 
//self.image_medium = resized_images['image_medium'] 
//self.image_small = resized_images['image_small'] 
//self.image = resized_images['image'] 
//if not self.image_medium: 
//self.image_medium = self.product_tmpl_id.image_medium 
//if not self.image_small: 
//self.image_small = self.product_tmpl_id.image_small 
//if not self.image: 
//self.image = self.product_tmpl_id.image 
})
pool.ProductProduct().Method().SetImage().DeclareMethod(
`SetImage` ,
func (){//def _set_image(self): 
//self._set_image_value(self.image) 
})
pool.ProductProduct().Method().SetImageMedium().DeclareMethod(
`SetImageMedium` ,
func (){//def _set_image_medium(self): 
//self._set_image_value(self.image_medium) 
})
pool.ProductProduct().Method().SetImageSmall().DeclareMethod(
`SetImageSmall` ,
func (){//def _set_image_small(self): 
//self._set_image_value(self.image_small) 
})
pool.ProductProduct().Method().SetImageValue().DeclareMethod(
`SetImageValue` ,
func (){//def _set_image_value(self, value): 
//image = tools.image_resize_image_big(value) 
//if self.product_tmpl_id.image: 
//self.image_variant = image 
//else: 
//self.product_tmpl_id.image = image 
})
pool.ProductProduct().Method().GetPricelistItems().DeclareMethod(
`GetPricelistItems` ,
func (){//def _get_pricelist_items(self): 
//self.pricelist_item_ids = self.env['product.pricelist.item'].search([ 
//'|', 
//('product_id', '=', self.id), 
//('product_tmpl_id', '=', self.product_tmpl_id.id)]).ids 
})
pool.ProductProduct().Method().CheckAttributeValueIds().DeclareMethod(
`CheckAttributeValueIds` ,
func (){//def _check_attribute_value_ids(self): 
//for product in self: 
//attributes = self.env['product.attribute'] 
//for value in product.attribute_value_ids: 
//if value.attribute_id in attributes: 
//raise ValidationError(_('Error! It is not allowed to choose more than one value for a given attribute.')) 
//attributes |= value.attribute_id 
//return True 
})
pool.ProductProduct().Method().OnchangeUom().DeclareMethod(
`OnchangeUom` ,
func (){//def _onchange_uom(self): 
//if self.uom_id and self.uom_po_id and self.uom_id.category_id != self.uom_po_id.category_id: 
//self.uom_po_id = self.uom_id 
})
pool.ProductProduct().Method().Create().DeclareMethod(
`Create` ,
func (){//def create(self, vals): 
//product = super(ProductProduct, self.with_context(create_product_product=True)).create(vals) 
//product._set_standard_price(vals.get('standard_price', 0.0)) 
//return product 
})
pool.ProductProduct().Method().Write().DeclareMethod(
`Write` ,
func (){//def write(self, values): 
//''' Store the standard price change in order to be able to retrieve the cost of a product for a given date''' 
//res = super(ProductProduct, self).write(values) 
//if 'standard_price' in values: 
//self._set_standard_price(values['standard_price']) 
//return res 
})
pool.ProductProduct().Method().Unlink().DeclareMethod(
`Unlink` ,
func (){//def unlink(self): 
//unlink_products = self.env['product.product'] 
//unlink_templates = self.env['product.template'] 
//for product in self: 
//# Check if product still exists, in case it has been unlinked by unlinking its template 
//if not product.exists(): 
//continue 
//# Check if the product is last product of this template 
//other_products = self.search([('product_tmpl_id', '=', product.product_tmpl_id.id), ('id', '!=', product.id)]) 
//if not other_products: 
//unlink_templates |= product.product_tmpl_id 
//unlink_products |= product 
//res = super(ProductProduct, unlink_products).unlink() 
//# delete templates after calling super, as deleting template could lead to deleting 
//# products due to ondelete='cascade' 
//unlink_templates.unlink() 
//return res 
})
pool.ProductProduct().Method().Copy().DeclareMethod(
`Copy` ,
func (){//def copy(self, default=None): 
//# TDE FIXME: clean context / variant brol 
//if default is None: 
//default = {} 
//if self._context.get('variant'): 
//# if we copy a variant or create one, we keep the same template 
//default['product_tmpl_id'] = self.product_tmpl_id.id 
//elif 'name' not in default: 
//default['name'] = self.name 
//return super(ProductProduct, self).copy(default=default) 
})
pool.ProductProduct().Method().Search().DeclareMethod(
`Search` ,
func (){//def search(self, args, offset=0, limit=None, order=None, count=False): 
//# TDE FIXME: strange 
//if self._context.get('search_default_categ_id'): 
//args.append((('categ_id', 'child_of', self._context['search_default_categ_id']))) 
//return super(ProductProduct, self).search(args, offset=offset, limit=limit, order=order, count=count) 
})
pool.ProductProduct().Method().NameGet().DeclareMethod(
`NameGet` ,
func (){//def name_get(self): 
})
pool.ProductProduct().Method().NameGet().DeclareMethod(
`NameGet` ,
func (){//def _name_get(d): 
//name = d.get('name', '') 
//code = self._context.get('display_default_code', True) and d.get('default_code', False) or False 
//if code: 
//name = '[%s] %s' % (code,name) 
//return (d['id'], name) 
//partner_id = self._context.get('partner_id') 
//if partner_id: 
//partner_ids = [partner_id, self.env['res.partner'].browse(partner_id).commercial_partner_id.id] 
//else: 
//partner_ids = [] 
//# all user don't have access to seller and partner 
//# check access and use superuser 
//self.check_access_rights("read") 
//self.check_access_rule("read") 
//result = [] 
//for product in self.sudo(): 
//# display only the attributes with multiple possible values on the template 
//variable_attributes = product.attribute_line_ids.filtered(lambda l: len(l.value_ids) > 1).mapped('attribute_id') 
//variant = product.attribute_value_ids._variant_name(variable_attributes) 
//name = variant and "%s (%s)" % (product.name, variant) or product.name 
//sellers = [] 
//if partner_ids: 
//sellers = [x for x in product.seller_ids if (x.name.id in partner_ids) and (x.product_id == product)] 
//if not sellers: 
//sellers = [x for x in product.seller_ids if (x.name.id in partner_ids) and not x.product_id] 
//if sellers: 
//for s in sellers: 
//seller_variant = s.product_name and ( 
//variant and "%s (%s)" % (s.product_name, variant) or s.product_name 
//) or False 
//mydict = { 
//'id': product.id, 
//'name': seller_variant or name, 
//'default_code': s.product_code or product.default_code, 
//} 
//temp = _name_get(mydict) 
//if temp not in result: 
//result.append(temp) 
//else: 
//mydict = { 
//'id': product.id, 
//'name': name, 
//'default_code': product.default_code, 
//} 
//result.append(_name_get(mydict)) 
//return result 
})
pool.ProductProduct().Method().NameSearch().DeclareMethod(
`NameSearch` ,
func (){//def name_search(self, name='', args=None, operator='ilike', limit=100): 
//if not args: 
//args = [] 
//if name: 
//positive_operators = ['=', 'ilike', '=ilike', 'like', '=like'] 
//products = self.env['product.product'] 
//if operator in positive_operators: 
//products = self.search([('default_code', '=', name)] + args, limit=limit) 
//if not products: 
//products = self.search([('barcode', '=', name)] + args, limit=limit) 
//if not products and operator not in expression.NEGATIVE_TERM_OPERATORS: 
//# Do not merge the 2 next lines into one single search, SQL search performance would be abysmal 
//# on a database with thousands of matching products, due to the huge merge+unique needed for the 
//# OR operator (and given the fact that the 'name' lookup results come from the ir.translation table 
//# Performing a quick memory merge of ids in Python will give much better performance 
//products = self.search(args + [('default_code', operator, name)], limit=limit) 
//if not limit or len(products) < limit: 
//# we may underrun the limit because of dupes in the results, that's fine 
//limit2 = (limit - len(products)) if limit else False 
//products += self.search(args + [('name', operator, name), ('id', 'not in', products.ids)], limit=limit2) 
//elif not products and operator in expression.NEGATIVE_TERM_OPERATORS: 
//products = self.search(args + ['&', ('default_code', operator, name), ('name', operator, name)], limit=limit) 
//if not products and operator in positive_operators: 
//ptrn = re.compile('(\[(.*?)\])') 
//res = ptrn.search(name) 
//if res: 
//products = self.search([('default_code', '=', res.group(2))] + args, limit=limit) 
//# still no results, partner in context: search on supplier info as last hope to find something 
//if not products and self._context.get('partner_id'): 
//suppliers = self.env['product.supplierinfo'].search([ 
//('name', '=', self._context.get('partner_id')), 
//'|', 
//('product_code', operator, name), 
//('product_name', operator, name)]) 
//if suppliers: 
//products = self.search([('product_tmpl_id.seller_ids', 'in', suppliers.ids)], limit=limit) 
//else: 
//products = self.search(args, limit=limit) 
//return products.name_get() 
})
pool.ProductProduct().Method().ViewHeaderGet().DeclareMethod(
`ViewHeaderGet` ,
func (){//def view_header_get(self, view_id, view_type): 
//res = super(ProductProduct, self).view_header_get(view_id, view_type) 
//if self._context.get('categ_id'): 
//return _('Products: ') + self.env['product.category'].browse(self._context['categ_id']).name 
//return res 
})
pool.ProductProduct().Method().OpenProductTemplate().DeclareMethod(
`OpenProductTemplate` ,
func (){//def open_product_template(self): 
//""" Utility method used to add an "Open Template" button in product views """ 
//self.ensure_one() 
//return {'type': 'ir.actions.act_window', 
//'res_model': 'product.template', 
//'view_mode': 'form', 
//'res_id': self.product_tmpl_id.id, 
//'target': 'new'} 
})
pool.ProductProduct().Method().SelectSeller().DeclareMethod(
`SelectSeller` ,
func (){//def _select_seller(self, partner_id=False, quantity=0.0, date=None, uom_id=False): 
//self.ensure_one() 
//if date is None: 
//date = fields.Date.today() 
//res = self.env['product.supplierinfo'] 
//for seller in self.seller_ids: 
//# Set quantity in UoM of seller 
//quantity_uom_seller = quantity 
//if quantity_uom_seller and uom_id and uom_id != seller.product_uom: 
//quantity_uom_seller = uom_id._compute_quantity(quantity_uom_seller, seller.product_uom) 
//if seller.date_start and seller.date_start > date: 
//continue 
//if seller.date_end and seller.date_end < date: 
//continue 
//if partner_id and seller.name not in [partner_id, partner_id.parent_id]: 
//continue 
//if quantity_uom_seller < seller.min_qty: 
//continue 
//if seller.product_id and seller.product_id != self: 
//continue 
//res |= seller 
//break 
//return res 
})
pool.ProductProduct().Method().PriceCompute().DeclareMethod(
`PriceCompute` ,
func (){//def price_compute(self, price_type, uom=False, currency=False, company=False): 
//# TDE FIXME: delegate to template or not ? fields are reencoded here ... 
//# compatibility about context keys used a bit everywhere in the code 
//if not uom and self._context.get('uom'): 
//uom = self.env['product.uom'].browse(self._context['uom']) 
//if not currency and self._context.get('currency'): 
//currency = self.env['res.currency'].browse(self._context['currency']) 
//products = self 
//if price_type == 'standard_price': 
//# standard_price field can only be seen by users in base.group_user 
//# Thus, in order to compute the sale price from the cost for users not in this group 
//# We fetch the standard price as the superuser 
//products = self.with_context(force_company=company and company.id or self._context.get('force_company', self.env.user.company_id.id)).sudo() 
//prices = dict.fromkeys(self.ids, 0.0) 
//for product in products: 
//prices[product.id] = product[price_type] or 0.0 
//if price_type == 'list_price': 
//prices[product.id] += product.price_extra 
//if uom: 
//prices[product.id] = product.uom_id._compute_price(prices[product.id], uom) 
//# Convert from current user company currency to asked one 
//# This is right cause a field cannot be in more than one currency 
//if currency: 
//prices[product.id] = product.currency_id.compute(prices[product.id], currency) 
//return prices 
// 
//# compatibility to remove after v10 - DEPRECATED 
})
pool.ProductProduct().Method().PriceGet().DeclareMethod(
`PriceGet` ,
func (){//def price_get(self, ptype='list_price'): 
//return self.price_compute(ptype) 
})
pool.ProductProduct().Method().SetStandardPrice().DeclareMethod(
`SetStandardPrice` ,
func (){//def _set_standard_price(self, value): 
//''' Store the standard price change in order to be able to retrieve the cost of a product for a given date''' 
//PriceHistory = self.env['product.price.history'] 
//for product in self: 
//PriceHistory.create({ 
//'product_id': product.id, 
//'cost': value, 
//'company_id': self._context.get('force_company', self.env.user.company_id.id), 
//}) 
})
pool.ProductProduct().Method().GetHistoryPrice().DeclareMethod(
`GetHistoryPrice` ,
func (){//def get_history_price(self, company_id, date=None): 
//history = self.env['product.price.history'].search([ 
//('company_id', '=', company_id), 
//('product_id', 'in', self.ids), 
//('datetime', '<=', date or fields.Datetime.now())], limit=1) 
})
pool.ProductProduct().Method().NeedProcurement().DeclareMethod(
`NeedProcurement` ,
func (){//def _need_procurement(self): 
//# When sale/product is installed alone, there is no need to create procurements. Only 
//# sale_stock and sale_service need procurements 
})


pool.ProductPackaging().DeclareModel()
pool.ProductPackaging().AddCharField("Name", models.StringFieldParams{String :"Packaging Type", Required: true})
pool.ProductPackaging().AddIntegerField("Sequence", models.SimpleFieldParams{String :"Sequence", Default: func(models.Environment, models.FieldMap) interface{} {return 1} ,Help :"The first in the sequence is the default one.)"})
pool.ProductPackaging().AddMany2OneField("ProductTmplId",models.ForeignKeyFieldParams{String :"Product" , RelationModel: pool.Product.Template()})
pool.ProductPackaging().AddFloatField("Qty", models.FloatFieldParams{String :"Quantity per Package" ,Help :"The total number of products you can have per pallet or box.)"})


pool.SuppliferInfo().DeclareModel()
pool.SuppliferInfo().AddMany2OneField("Name",models.ForeignKeyFieldParams{String :"Vendor" , RelationModel: pool.Res.Partner()})
pool.SuppliferInfo().AddCharField("ProductName", models.StringFieldParams{String :"Vendor Product Name" ,Help :"This vendor's product name will be used when printing a request for quotation. Keep empty to use the internal one.)"})
pool.SuppliferInfo().AddCharField("ProductCode", models.StringFieldParams{String :"Vendor Product Code" ,Help :"This vendor's product code will be used when printing a request for quotation. Keep empty to use the internal one.)"})
pool.SuppliferInfo().AddIntegerField("Sequence", models.SimpleFieldParams{String :"Sequence", Default: func(models.Environment, models.FieldMap) interface{} {return 1} ,Help :"Assigns the priority to the list of product vendor.)"})
pool.SuppliferInfo().AddMany2OneField("ProductUom",models.ForeignKeyFieldParams{String :"Vendor Unit of Measure" , RelationModel: pool.Product.Uom(), Related: "ProductTmplId.UomPoId" , Help :"This comes from the product form.)"})
pool.SuppliferInfo().Fields().ProductUom().RevokeAccess(security.GroupEveryone, security.Write)
pool.SuppliferInfo().AddFloatField("MinQty", models.FloatFieldParams{String :"Minimal Quantity", Default: func(models.Environment, models.FieldMap) interface{} {return 0.0}, Required: true ,Help :"The minimal quantity to purchase from this vendor, expressed in the vendor Product Unit of Measure if not any, in the default unit of measure of the product otherwise.)"})
pool.SuppliferInfo().AddFloatField("Price", models.FloatFieldParams{String :"Price", Default: func(models.Environment, models.FieldMap) interface{} {return 0.0}})
pool.SuppliferInfo().AddMany2OneField("CompanyId",models.ForeignKeyFieldParams{String :"Company" , RelationModel: pool.Res.Company(), Default: func(models.Environment, models.FieldMap) interface{} {return lambda self: self.env.user.company_id.id}, Index: true})
pool.SuppliferInfo().AddMany2OneField("CurrencyId",models.ForeignKeyFieldParams{String :"Currency" , RelationModel: pool.Res.Currency(), Default: func(models.Environment, models.FieldMap) interface{} {return lambda self: self.env.user.company_id.currency_id.id}, Required: true})
pool.SuppliferInfo().AddDateField("DateStart", models.SimpleFieldParams{String :"Start Date" ,Help :"Start date for this vendor price)"})
pool.SuppliferInfo().AddDateField("DateEnd", models.SimpleFieldParams{String :"End Date" ,Help :"End date for this vendor price)"})
pool.SuppliferInfo().AddMany2OneField("ProductId",models.ForeignKeyFieldParams{String :"Product Variant" , RelationModel: pool.Product.Product() , Help :"When this field is filled in, the vendor data will only apply to the variant.)"})
pool.SuppliferInfo().AddMany2OneField("ProductTmplId",models.ForeignKeyFieldParams{String :"Product Template" , RelationModel: pool.Product.Template(), Index: true, OnDelete : models.Cascade})
pool.SuppliferInfo().AddIntegerField("Delay", models.SimpleFieldParams{String :"Delivery Lead Time", Default: func(models.Environment, models.FieldMap) interface{} {return 1}, Required: true ,Help :"Lead time in days between the confirmation of the purchase order and the receipt of the products in your warehouse. Used by the scheduler for automatic computation of the purchase order planning.)"})
 
 }