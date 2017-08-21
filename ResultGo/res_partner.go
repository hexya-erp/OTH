package res 

  import (

) 

 func init() { 

 

pool.Partner().DeclareModel()
pool.Partner().AddMany2OneField("PropertyProductPricelist",models.ForeignKeyFieldParams{String :"Sale Pricelist" , RelationModel: pool.ProductPricelist()})
pool.Partner().Methods().ComputeProductPricelist().DeclareMethod(
`ComputeProductPricelist` ,
func (rs pool.PartnerSet){
  //@api.depends('country_id')
  /*def _compute_product_pricelist(self):
        for p in self:
            if not isinstance(p.id, models.NewId):  # if not onchange
                p.property_product_pricelist = self.env['product.pricelist']._get_partner_pricelist(p.id)

    */})
pool.Partner().Methods().InverseProductPricelist().DeclareMethod(
`InverseProductPricelist` ,
func (rs pool.PartnerSet){
  //@api.one
  /*def _inverse_product_pricelist(self):
        pls = self.env['product.pricelist'].search(
            [('country_group_ids.country_ids.code', '=', self.country_id and self.country_id.code or False)],
            limit=1
        )
        default_for_country = pls and pls[0]
        actual = self.env['ir.property'].get('property_product_pricelist', 'res.partner', 'res.partner,%s' % self.id)

        # update at each change country, and so erase old pricelist
        if self.property_product_pricelist or (actual and default_for_country and default_for_country.id != actual.id):
            # keep the company of the current user before sudo
            self.env['ir.property'].with_context(force_company=self.env.user.company_id.id).sudo().set_multi(
                'property_product_pricelist',
                self._name,
                {self.id: self.property_product_pricelist or default_for_country.id},
                default_value=default_for_country.id
            )

    */})
pool.Partner().Methods().CommercialFields().DeclareMethod(
`CommercialFields` ,
func (rs pool.PartnerSet){
  /*def _commercial_fields(self):
        return super(Partner, self)._commercial_fields() + ['property_product_pricelist']
*/})
 
 }