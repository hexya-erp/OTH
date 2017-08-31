package product 

  import (
"github.com/hexya-erp/hexya/hexya/models/types"

) 

 func init() { 

 

pool.ProductCategory().DeclareModel()
pool.ProductCategory().AddCharField("Name", models.StringFieldParams{String :"Name", Index: true, Required : true, Translate: true})
pool.ProductCategory().AddMany2OneField("Parent",models.ForeignKeyFieldParams{String :"Parent Category" , RelationModel: pool.ProductCategory() , JSON : "parent_id", Index: true, OnDelete : models.Cascade})
pool.ProductCategory().AddOne2ManyField("Child", models.ReverseFieldParams{String :"Child Categories" ,RelationModel : pool.ProductCategory() ,ReverseFK : "Parent" , JSON : "child_id"})
pool.ProductCategory().AddSelectionField("Type", models.SelectionFieldParams{String :"Category Type", Selection : types.Selection{
"view" : "View",
"normal" : "Normal",
}, Default: func(models.Environment, models.FieldMap) interface{} {return "normal"}, Help : "A category of the view type is a virtual category that can be used as the parent of another category to create a hierarchical structure."})
pool.ProductCategory().AddIntegerField("ParentLeft", models.SimpleFieldParams{String :"Left Parent", Index: true})
pool.ProductCategory().AddIntegerField("ParentRight", models.SimpleFieldParams{String :"Right Parent", Index: true})
pool.ProductCategory().AddIntegerField("ProductCount", models.SimpleFieldParams{String :"# Products" ,Help :"The number of products under this category "})
pool.ProductCategory().Methods().ComputeProductCount().DeclareMethod(
`ComputeProductCount` ,
func (rs pool.ProductCategorySet){
  /*def _compute_product_count(self):
        read_group_res = self.env['product.template'].read_group([('categ_id', 'in', self.ids)], ['categ_id'], ['categ_id'])
        group_data = dict((data['categ_id'][0], data['categ_id_count']) for data in read_group_res)
        for categ in self:
            categ.product_count = group_data.get(categ.id, 0)

    */})
pool.ProductCategory().Methods().CheckCategoryRecursion().DeclareMethod(
`CheckCategoryRecursion` ,
func (rs pool.ProductCategorySet){
  //@api.constrains('parent_id')
  /*def _check_category_recursion(self):
        if not self._check_recursion():
            raise ValidationError(_('Error ! You cannot create recursive categories.'))
        return True

    */})
pool.ProductCategory().Methods().NameGet().DeclareMethod(
`NameGet` ,
func (rs pool.ProductCategorySet){
  //@api.multi
  /*def name_get(self):
        */})
pool.ProductCategory().Methods().GetNames().DeclareMethod(
`GetNames` ,
func (rs pool.ProductCategorySet){
  /*def get_names(cat):
            """ Return the list [cat.name, cat.parent_id.name, ...] """
            res = []
            while cat:
                res.append(cat.name)
                cat = cat.parent_id
            return res

        return [(cat.id, " / ".join(reversed(get_names(cat)))) for cat in self]

    */})
pool.ProductCategory().Methods().NameSearch().DeclareMethod(
`NameSearch` ,
func (rs pool.ProductCategorySet , args struct{Args interface{}
Operator interface{}
Limit interface{}
}){
  //@api.model
  /*def name_search(self, name, args=None, operator='ilike', limit=100):
        if not args:
            args = []
        if name:
            # Be sure name_search is symetric to name_get
            category_names = name.split(' / ')
            parents = list(category_names)
            child = parents.pop()
            domain = [('name', operator, child)]
            if parents:
                names_ids = self.name_search(' / '.join(parents), args=args, operator='ilike', limit=limit)
                category_ids = [name_id[0] for name_id in names_ids]
                if operator in expression.NEGATIVE_TERM_OPERATORS:
                    categories = self.search([('id', 'not in', category_ids)])
                    domain = expression.OR([[('parent_id', 'in', categories.ids)], domain])
                else:
                    domain = expression.AND([[('parent_id', 'in', category_ids)], domain])
                for i in range(1, len(category_names)):
                    domain = [[('name', operator, ' / '.join(category_names[-1 - i:]))], domain]
                    if operator in expression.NEGATIVE_TERM_OPERATORS:
                        domain = expression.AND(domain)
                    else:
                        domain = expression.OR(domain)
            categories = self.search(expression.AND([domain, args]), limit=limit)
        else:
            categories = self.search(args, limit=limit)
        return categories.name_get()


*/})


pool.ProductPriceHistory().DeclareModel()
pool.ProductPriceHistory().Methods().GetDefaultCompanyId().DeclareMethod(
`GetDefaultCompanyId` ,
func (rs pool.ProductPriceHistorySet){
  /*def _get_default_company_id(self):
        return self._context.get('force_company', self.env.user.company_id.id)

    company_id = */})
pool.ProductPriceHistory().AddMany2OneField("Company",models.ForeignKeyFieldParams{String :"Company" , RelationModel: pool.Company() , JSON : "company_id", Default : func(models.Environment, models.FieldMap) interface{}{
/*_get_default_company_id(self):
        return self._context.get('force_company', self.env.user.company_id.id)

    company_id = */
return 0}, Required : true})
pool.ProductPriceHistory().AddMany2OneField("Product",models.ForeignKeyFieldParams{String :"Product" , RelationModel: pool.ProductProduct() , JSON : "product_id", OnDelete : models.Cascade, Required : true})
pool.ProductPriceHistory().AddDateTimeField("Datetime", models.SimpleFieldParams{})
pool.ProductPriceHistory().AddFloatField("Cost", models.FloatFieldParams{String :"Cost"})


pool.ProductProduct().DeclareModel()
pool.ProductProduct().AddFloatField("Price", models.FloatFieldParams{String :"Price"})
pool.ProductProduct().AddFloatField("PriceExtra", models.FloatFieldParams{String :"Variant Price Extra"})
pool.ProductProduct().AddFloatField("LstPrice", models.FloatFieldParams{String :"Sale Price"})
pool.ProductProduct().AddCharField("DefaultCode", models.StringFieldParams{String :"Internal Reference", Index: true})
pool.ProductProduct().AddCharField("Code", models.StringFieldParams{String :"Internal Reference"})
pool.ProductProduct().AddCharField("PartnerRef", models.StringFieldParams{String :"Customer Ref"})
pool.ProductProduct().AddBooleanField("Active", models.SimpleFieldParams{String :"Active", Default: func(models.Environment, models.FieldMap) interface{} {return true} ,Help :"If unchecked, it will allow you to hide the product without removing it."})
pool.ProductProduct().AddMany2OneField("ProductTmpl",models.ForeignKeyFieldParams{String :"Product Template" , RelationModel: pool.ProductTemplate() , JSON : "product_tmpl_id", Index: true, OnDelete : models.Cascade, Required : true})
pool.ProductProduct().AddCharField("Barcode", models.StringFieldParams{String :"Barcode", NoCopy: true ,Help :"International Article Number used for product identification."})
pool.ProductProduct().AddMany2ManyField("AttributeValues", models.Many2ManyFieldParams{String :"Attributes" , RelationModel: pool.ProductAttributeValue() , JSON : "attribute_value_ids"})
pool.ProductProduct().AddBinaryField("ImageVariant", models.SimpleFieldParams{String :"Variant Image" ,Help :"This field holds the image used as image for the product variant, limited to 1024x1024px."})
pool.ProductProduct().AddBinaryField("ImageSmall", models.SimpleFieldParams{String :"Small-sized image" ,Help :"Image of the product variant "})
pool.ProductProduct().AddBinaryField("Image", models.SimpleFieldParams{String :"Big-sized image" ,Help :"Image of the product variant "})
pool.ProductProduct().AddBinaryField("ImageMedium", models.SimpleFieldParams{String :"Medium-sized image" ,Help :"Image of the product variant "})
pool.ProductProduct().AddFloatField("StandardPrice", models.FloatFieldParams{String :"Cost"})
pool.ProductProduct().AddFloatField("Volume", models.FloatFieldParams{String :"Volume" ,Help :"The volume in m3."})
pool.ProductProduct().AddFloatField("Weight", models.FloatFieldParams{String :"Weight"})
pool.ProductProduct().AddMany2ManyField("PricelistItems", models.Many2ManyFieldParams{String :"Pricelist Items" , RelationModel: pool.ProductPricelistItem() , JSON : "pricelist_item_ids"})
pool.ProductProduct().AddSQLConstraint("BarcodeUniq" , "unique(barcode)" , "A barcode can only be assigned to one product !")
pool.ProductProduct().Methods().ComputeProductPrice().DeclareMethod(
`ComputeProductPrice` ,
func (rs pool.ProductProductSet){
  /*def _compute_product_price(self):
        prices = {}
        pricelist_id_or_name = self._context.get('pricelist')
        if pricelist_id_or_name:
            pricelist = None
            partner = self._context.get('partner', False)
            quantity = self._context.get('quantity', 1.0)

            # Support context pricelists specified as display_name or ID for compatibility
            if isinstance(pricelist_id_or_name, basestring):
                pricelist_name_search = self.env['product.pricelist'].name_search(pricelist_id_or_name, operator='=', limit=1)
                if pricelist_name_search:
                    pricelist = self.env['product.pricelist'].browse([pricelist_name_search[0][0]])
            elif isinstance(pricelist_id_or_name, (int, long)):
                pricelist = self.env['product.pricelist'].browse(pricelist_id_or_name)

            if pricelist:
                quantities = [quantity] * len(self)
                partners = [partner] * len(self)
                prices = pricelist.get_products_price(self, quantities, partners)

        for product in self:
            product.price = prices.get(product.id, 0.0)

    */})
pool.ProductProduct().Methods().InverseProductPrice().DeclareMethod(
`InverseProductPrice` ,
func (rs pool.ProductProductSet){
  /*def _set_product_price(self):
        for product in self:
            if self._context.get('uom'):
                value = self.env['product.uom'].browse(self._context['uom'])._compute_price(product.price, product.uom_id)
            else:
                value = product.price
            value -= product.price_extra
            product.write({'list_price': value})

    */})
pool.ProductProduct().Methods().InverseProductLstPrice().DeclareMethod(
`InverseProductLstPrice` ,
func (rs pool.ProductProductSet){
  /*def _set_product_lst_price(self):
        for product in self:
            if self._context.get('uom'):
                value = self.env['product.uom'].browse(self._context['uom'])._compute_price(product.lst_price, product.uom_id)
            else:
                value = product.lst_price
            value -= product.price_extra
            product.write({'list_price': value})

    */})
pool.ProductProduct().Methods().ComputeProductPriceExtra().DeclareMethod(
`ComputeProductPriceExtra` ,
func (rs pool.ProductProductSet){
  //@api.depends('attribute_value_ids.price_ids.price_extra','attribute_value_ids.price_ids.product_tmpl_id')
  /*def _compute_product_price_extra(self):
        # TDE FIXME: do a real multi and optimize a bit ?
        for product in self:
            price_extra = 0.0
            for attribute_price in product.mapped('attribute_value_ids.price_ids'):
                if attribute_price.product_tmpl_id == product.product_tmpl_id:
                    price_extra += attribute_price.price_extra
            product.price_extra = price_extra

    */})
pool.ProductProduct().Methods().ComputeProductLstPrice().DeclareMethod(
`ComputeProductLstPrice` ,
func (rs pool.ProductProductSet){
  //@api.depends('list_price','price_extra')
  /*def _compute_product_lst_price(self):
        to_uom = None
        if 'uom' in self._context:
            to_uom = self.env['product.uom'].browse([self._context['uom']])

        for product in self:
            if to_uom:
                list_price = product.uom_id._compute_price(product.list_price, to_uom)
            else:
                list_price = product.list_price
            product.lst_price = list_price + product.price_extra

    */})
pool.ProductProduct().Methods().ComputeProductCode().DeclareMethod(
`ComputeProductCode` ,
func (rs pool.ProductProductSet){
  //@api.one
  /*def _compute_product_code(self):
        for supplier_info in self.seller_ids:
            if supplier_info.name.id == self._context.get('partner_id'):
                self.code = supplier_info.product_code or self.default_code
        else:
            self.code = self.default_code

    */})
pool.ProductProduct().Methods().ComputePartnerRef().DeclareMethod(
`ComputePartnerRef` ,
func (rs pool.ProductProductSet){
  //@api.one
  /*def _compute_partner_ref(self):
        for supplier_info in self.seller_ids:
            if supplier_info.name.id == self._context.get('partner_id'):
                product_name = supplier_info.product_name or self.default_code
        else:
            product_name = self.name
        self.partner_ref = '%s%s' % (self.code and '[%s] ' % self.code or '', product_name)

    */})
pool.ProductProduct().Methods().ComputeImages().DeclareMethod(
`ComputeImages` ,
func (rs pool.ProductProductSet){
  //@api.depends('image_variant','product_tmpl_id.image')
  /*def _compute_images(self):
        if self._context.get('bin_size'):
            self.image_medium = self.image_variant
            self.image_small = self.image_variant
            self.image = self.image_variant
        else:
            resized_images = tools.image_get_resized_images(self.image_variant, return_big=True, avoid_resize_medium=True)
            self.image_medium = resized_images['image_medium']
            self.image_small = resized_images['image_small']
            self.image = resized_images['image']
        if not self.image_medium:
            self.image_medium = self.product_tmpl_id.image_medium
        if not self.image_small:
            self.image_small = self.product_tmpl_id.image_small
        if not self.image:
            self.image = self.product_tmpl_id.image

    */})
pool.ProductProduct().Methods().InverseImage().DeclareMethod(
`InverseImage` ,
func (rs pool.ProductProductSet){
  //@api.one
  /*def _set_image(self):
        self._set_image_value(self.image)

    */})
pool.ProductProduct().Methods().InverseImageMedium().DeclareMethod(
`InverseImageMedium` ,
func (rs pool.ProductProductSet){
  //@api.one
  /*def _set_image_medium(self):
        self._set_image_value(self.image_medium)

    */})
pool.ProductProduct().Methods().InverseImageSmall().DeclareMethod(
`InverseImageSmall` ,
func (rs pool.ProductProductSet){
  //@api.one
  /*def _set_image_small(self):
        self._set_image_value(self.image_small)

    */})
pool.ProductProduct().Methods().InverseImageValue().DeclareMethod(
`InverseImageValue` ,
func (rs pool.ProductProductSet , args struct{Value interface{}
}){
  //@api.one
  /*def _set_image_value(self, value):
        image = tools.image_resize_image_big(value)
        if self.product_tmpl_id.image:
            self.image_variant = image
        else:
            self.product_tmpl_id.image = image

    */})
pool.ProductProduct().Methods().GetPricelistItems().DeclareMethod(
`GetPricelistItems` ,
func (rs pool.ProductProductSet){
  //@api.one
  /*def _get_pricelist_items(self):
        self.pricelist_item_ids = self.env['product.pricelist.item'].search([
            '|',
            ('product_id', '=', self.id),
            ('product_tmpl_id', '=', self.product_tmpl_id.id)]).ids

    */})
pool.ProductProduct().Methods().CheckAttributeValueIds().DeclareMethod(
`CheckAttributeValueIds` ,
func (rs pool.ProductProductSet){
  //@api.constrains('attribute_value_ids')
  /*def _check_attribute_value_ids(self):
        for product in self:
            attributes = self.env['product.attribute']
            for value in product.attribute_value_ids:
                if value.attribute_id in attributes:
                    raise ValidationError(_('Error! It is not allowed to choose more than one value for a given attribute.'))
                attributes |= value.attribute_id
        return True

    */})
pool.ProductProduct().Methods().OnchangeUom().DeclareMethod(
`OnchangeUom` ,
func (rs pool.ProductProductSet){
  //@api.onchange('uom_id','uom_po_id')
  /*def _onchange_uom(self):
        if self.uom_id and self.uom_po_id and self.uom_id.category_id != self.uom_po_id.category_id:
            self.uom_po_id = self.uom_id

    */})
pool.ProductProduct().Methods().Create().DeclareMethod(
`Create` ,
func (rs pool.ProductProductSet , args struct{}){
  //@api.model
  /*def create(self, vals):
        product = super(ProductProduct, self.with_context(create_product_product=True)).create(vals)
        product._set_standard_price(vals.get('standard_price', 0.0))
        return product

    */})
pool.ProductProduct().Methods().Write().DeclareMethod(
`Write` ,
func (rs pool.ProductProductSet , args struct{Values interface{}
}){
  //@api.multi
  /*def write(self, values):
        ''' Store the standard price change in order to be able to retrieve the cost of a product for a given date'''
        res = super(ProductProduct, self).write(values)
        if 'standard_price' in values:
            self._set_standard_price(values['standard_price'])
        return res

    */})
pool.ProductProduct().Methods().Unlink().DeclareMethod(
`Unlink` ,
func (rs pool.ProductProductSet){
  //@api.multi
  /*def unlink(self):
        unlink_products = self.env['product.product']
        unlink_templates = self.env['product.template']
        for product in self:
            # Check if product still exists, in case it has been unlinked by unlinking its template
            if not product.exists():
                continue
            # Check if the product is last product of this template
            other_products = self.search([('product_tmpl_id', '=', product.product_tmpl_id.id), ('id', '!=', product.id)])
            if not other_products:
                unlink_templates |= product.product_tmpl_id
            unlink_products |= product
        res = super(ProductProduct, unlink_products).unlink()
        # delete templates after calling super, as deleting template could lead to deleting
        # products due to ondelete='cascade'
        unlink_templates.unlink()
        return res

    */})
pool.ProductProduct().Methods().Copy().DeclareMethod(
`Copy` ,
func (rs pool.ProductProductSet , args struct{Default interface{}
}){
  //@api.multi
  /*def copy(self, default=None):
        # TDE FIXME: clean context / variant brol
        if default is None:
            default = {}
        if self._context.get('variant'):
            # if we copy a variant or create one, we keep the same template
            default['product_tmpl_id'] = self.product_tmpl_id.id
        elif 'name' not in default:
            default['name'] = self.name

        return super(ProductProduct, self).copy(default=default)

    */})
pool.ProductProduct().Methods().Search().DeclareMethod(
`Search` ,
func (rs pool.ProductProductSet , args struct{Offset interface{}
Limit interface{}
Order interface{}
Count interface{}
}){
  //@api.model
  /*def search(self, args, offset=0, limit=None, order=None, count=False):
        # TDE FIXME: strange
        if self._context.get('search_default_categ_id'):
            args.append((('categ_id', 'child_of', self._context['search_default_categ_id'])))
        return super(ProductProduct, self).search(args, offset=offset, limit=limit, order=order, count=count)

    */})
pool.ProductProduct().Methods().NameGet().DeclareMethod(
`NameGet` ,
func (rs pool.ProductProductSet){
  //@api.multi
  /*def name_get(self):
        */})
pool.ProductProduct().Methods().NameGet().DeclareMethod(
`NameGet` ,
func (rs pool.ProductProductSet){
  /*def _name_get(d):
            name = d.get('name', '')
            code = self._context.get('display_default_code', True) and d.get('default_code', False) or False
            if code:
                name = '[%s] %s' % (code,name)
            return (d['id'], name)

        partner_id = self._context.get('partner_id')
        if partner_id:
            partner_ids = [partner_id, self.env['res.partner'].browse(partner_id).commercial_partner_id.id]
        else:
            partner_ids = []

        # all user don't have access to seller and partner
        # check access and use superuser
        self.check_access_rights("read")
        self.check_access_rule("read")

        result = []
        for product in self.sudo():
            # display only the attributes with multiple possible values on the template
            variable_attributes = product.attribute_line_ids.filtered(lambda l: len(l.value_ids) > 1).mapped('attribute_id')
            variant = product.attribute_value_ids._variant_name(variable_attributes)

            name = variant and "%s (%s)" % (product.name, variant) or product.name
            sellers = []
            if partner_ids:
                sellers = [x for x in product.seller_ids if (x.name.id in partner_ids) and (x.product_id == product)]
                if not sellers:
                    sellers = [x for x in product.seller_ids if (x.name.id in partner_ids) and not x.product_id]
            if sellers:
                for s in sellers:
                    seller_variant = s.product_name and (
                        variant and "%s (%s)" % (s.product_name, variant) or s.product_name
                        ) or False
                    mydict = {
                              'id': product.id,
                              'name': seller_variant or name,
                              'default_code': s.product_code or product.default_code,
                              }
                    temp = _name_get(mydict)
                    if temp not in result:
                        result.append(temp)
            else:
                mydict = {
                          'id': product.id,
                          'name': name,
                          'default_code': product.default_code,
                          }
                result.append(_name_get(mydict))
        return result

    */})
pool.ProductProduct().Methods().NameSearch().DeclareMethod(
`NameSearch` ,
func (rs pool.ProductProductSet , args struct{Args interface{}
Operator interface{}
Limit interface{}
}){
  //@api.model
  /*def name_search(self, name, args=None, operator='ilike', limit=100):
        if not args:
            args = []
        if name:
            # Be sure name_search is symetric to name_get
            category_names = name.split(' / ')
            parents = list(category_names)
            child = parents.pop()
            domain = [('name', operator, child)]
            if parents:
                names_ids = self.name_search(' / '.join(parents), args=args, operator='ilike', limit=limit)
                category_ids = [name_id[0] for name_id in names_ids]
                if operator in expression.NEGATIVE_TERM_OPERATORS:
                    categories = self.search([('id', 'not in', category_ids)])
                    domain = expression.OR([[('parent_id', 'in', categories.ids)], domain])
                else:
                    domain = expression.AND([[('parent_id', 'in', category_ids)], domain])
                for i in range(1, len(category_names)):
                    domain = [[('name', operator, ' / '.join(category_names[-1 - i:]))], domain]
                    if operator in expression.NEGATIVE_TERM_OPERATORS:
                        domain = expression.AND(domain)
                    else:
                        domain = expression.OR(domain)
            categories = self.search(expression.AND([domain, args]), limit=limit)
        else:
            categories = self.search(args, limit=limit)
        return categories.name_get()


*/})
pool.ProductProduct().Methods().ViewHeaderGet().DeclareMethod(
`ViewHeaderGet` ,
func (rs pool.ProductProductSet , args struct{ViewId interface{}
ViewType interface{}
}){
  //@api.model
  /*def view_header_get(self, view_id, view_type):
        res = super(ProductProduct, self).view_header_get(view_id, view_type)
        if self._context.get('categ_id'):
            return _('Products: ') + self.env['product.category'].browse(self._context['categ_id']).name
        return res

    */})
pool.ProductProduct().Methods().OpenProductTemplate().DeclareMethod(
`OpenProductTemplate` ,
func (rs pool.ProductProductSet){
  //@api.multi
  /*def open_product_template(self):
        """ Utility method used to add an "Open Template" button in product views """
        self.ensure_one()
        return {'type': 'ir.actions.act_window',
                'res_model': 'product.template',
                'view_mode': 'form',
                'res_id': self.product_tmpl_id.id,
                'target': 'new'}

    */})
pool.ProductProduct().Methods().SelectSeller().DeclareMethod(
`SelectSeller` ,
func (rs pool.ProductProductSet , args struct{PartnerId interface{}
Quantity interface{}
Date interface{}
UomId interface{}
}){
  //@api.multi
  /*def _select_seller(self, partner_id=False, quantity=0.0, date=None, uom_id=False):
        self.ensure_one()
        if date is None:
            date = */})
pool.ProductProduct().Methods().PriceCompute().DeclareMethod(
`PriceCompute` ,
func (rs pool.ProductProductSet , args struct{PriceType interface{}
Uom interface{}
Currency interface{}
Company interface{}
}){
  //@api.multi
  /*def price_compute(self, price_type, uom=False, currency=False, company=False):
        # TDE FIXME: delegate to template or not ? fields are reencoded here ...
        # compatibility about context keys used a bit everywhere in the code
        if not uom and self._context.get('uom'):
            uom = self.env['product.uom'].browse(self._context['uom'])
        if not currency and self._context.get('currency'):
            currency = self.env['res.currency'].browse(self._context['currency'])

        products = self
        if price_type == 'standard_price':
            # standard_price field can only be seen by users in base.group_user
            # Thus, in order to compute the sale price from the cost for users not in this group
            # We fetch the standard price as the superuser
            products = self.with_context(force_company=company and company.id or self._context.get('force_company', self.env.user.company_id.id)).sudo()

        prices = dict.fromkeys(self.ids, 0.0)
        for product in products:
            prices[product.id] = product[price_type] or 0.0
            if price_type == 'list_price':
                prices[product.id] += product.price_extra

            if uom:
                prices[product.id] = product.uom_id._compute_price(prices[product.id], uom)

            # Convert from current user company currency to asked one
            # This is right cause a field cannot be in more than one currency
            if currency:
                prices[product.id] = product.currency_id.compute(prices[product.id], currency)

        return prices


    # compatibility to remove after v10 - DEPRECATED
    */})
pool.ProductProduct().Methods().PriceGet().DeclareMethod(
`PriceGet` ,
func (rs pool.ProductProductSet , args struct{Ptype interface{}
}){
  //@api.multi
  /*def price_get(self, ptype='list_price'):
        return self.price_compute(ptype)

    */})
pool.ProductProduct().Methods().InverseStandardPrice().DeclareMethod(
`InverseStandardPrice` ,
func (rs pool.ProductProductSet , args struct{Value interface{}
}){
  //@api.multi
  /*def _set_standard_price(self, value):
        ''' Store the standard price change in order to be able to retrieve the cost of a product for a given date'''
        PriceHistory = self.env['product.price.history']
        for product in self:
            PriceHistory.create({
                'product_id': product.id,
                'cost': value,
                'company_id': self._context.get('force_company', self.env.user.company_id.id),
            })

    */})
pool.ProductProduct().Methods().GetHistoryPrice().DeclareMethod(
`GetHistoryPrice` ,
func (rs pool.ProductProductSet , args struct{CompanyId interface{}
Date interface{}
}){
  //@api.multi
  /*def get_history_price(self, company_id, date=None):
        history = self.env['product.price.history'].search([
            ('company_id', '=', company_id),
            ('product_id', 'in', self.ids),
            ('datetime', '<=', date or */})
pool.ProductProduct().Methods().NeedProcurement().DeclareMethod(
`NeedProcurement` ,
func (rs pool.ProductProductSet){
  /*def _need_procurement(self):
        # When sale/product is installed alone, there is no need to create procurements. Only
        # sale_stock and sale_service need procurements
        return False


*/})


pool.ProductPackaging().DeclareModel()
pool.ProductPackaging().AddCharField("Name", models.StringFieldParams{String :"Packaging Type", Required : true})
pool.ProductPackaging().AddIntegerField("Sequence", models.SimpleFieldParams{String :"Sequence", Default: func(models.Environment, models.FieldMap) interface{} {return 1} ,Help :"The first in the sequence is the default one."})
pool.ProductPackaging().AddMany2OneField("ProductTmpl",models.ForeignKeyFieldParams{String :"Product" , RelationModel: pool.ProductTemplate() , JSON : "product_tmpl_id"})
pool.ProductPackaging().AddFloatField("Qty", models.FloatFieldParams{String :"Quantity per Package" ,Help :"The total number of products you can have per pallet or box."})


pool.ProductSupplierinfo().DeclareModel()
pool.ProductSupplierinfo().AddMany2OneField("Name",models.ForeignKeyFieldParams{String :"Vendor" , RelationModel: pool.Partner()})
pool.ProductSupplierinfo().AddCharField("ProductName", models.StringFieldParams{String :"Vendor Product Name" ,Help :"This vendor's product name will be used when printing a request for quotation. Keep empty to use the internal one."})
pool.ProductSupplierinfo().AddCharField("ProductCode", models.StringFieldParams{String :"Vendor Product Code" ,Help :"This vendor's product code will be used when printing a request for quotation. Keep empty to use the internal one."})
pool.ProductSupplierinfo().AddIntegerField("Sequence", models.SimpleFieldParams{String :"Sequence", Default: func(models.Environment, models.FieldMap) interface{} {return 1} ,Help :"Assigns the priority to the list of product vendor."})
pool.ProductSupplierinfo().AddMany2OneField("ProductUom",models.ForeignKeyFieldParams{String :"Vendor Unit of Measure" , RelationModel: pool.ProductUom(), Related: "ProductTmpl.UomPo" , Help :"This comes from the product form."})
pool.ProductSupplierinfo().Fields().ProductUom().RevokeAccess(security.GroupEveryone, security.Write)
pool.ProductSupplierinfo().AddFloatField("MinQty", models.FloatFieldParams{String :"Minimal Quantity", Default: func(models.Environment, models.FieldMap) interface{} {return 0.0}, Required : true ,Help :"The minimal quantity to purchase from this vendor, expressed in the vendor Product Unit of Measure if not any, in the default unit of measure of the product otherwise."})
pool.ProductSupplierinfo().AddFloatField("Price", models.FloatFieldParams{String :"Price", Default: func(models.Environment, models.FieldMap) interface{} {return 0.0}})
pool.ProductSupplierinfo().AddMany2OneField("Company",models.ForeignKeyFieldParams{String :"Company" , RelationModel: pool.Company() , JSON : "company_id", Default : func(models.Environment, models.FieldMap) interface{}{
/*lambda self: self.env.user.company_id.id*/
return 0}, Index: true})
pool.ProductSupplierinfo().AddMany2OneField("Currency",models.ForeignKeyFieldParams{String :"Currency" , RelationModel: pool.Currency() , JSON : "currency_id", Default : func(models.Environment, models.FieldMap) interface{}{
/*lambda self: self.env.user.company_id.currency_id.id*/
return 0}, Required : true})
pool.ProductSupplierinfo().AddDateField("DateStart", models.SimpleFieldParams{String :"Start Date" ,Help :"Start date for this vendor price"})
pool.ProductSupplierinfo().AddDateField("DateEnd", models.SimpleFieldParams{String :"End Date" ,Help :"End date for this vendor price"})
pool.ProductSupplierinfo().AddMany2OneField("Product",models.ForeignKeyFieldParams{String :"Product Variant" , RelationModel: pool.ProductProduct() , JSON : "product_id" , Help :"When this field is filled in, the vendor data will only apply to the variant."})
pool.ProductSupplierinfo().AddMany2OneField("ProductTmpl",models.ForeignKeyFieldParams{String :"Product Template" , RelationModel: pool.ProductTemplate() , JSON : "product_tmpl_id", Index: true, OnDelete : models.Cascade})
pool.ProductSupplierinfo().AddIntegerField("Delay", models.SimpleFieldParams{String :"Delivery Lead Time", Default: func(models.Environment, models.FieldMap) interface{} {return 1}, Required : true ,Help :"Lead time in days between the confirmation of the purchase order and the receipt of the products in your warehouse. Used by the scheduler for automatic computation of the purchase order planning."})
 
 }