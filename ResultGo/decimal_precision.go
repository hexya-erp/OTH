package  

  import (

) 

 func init() { 

 

pool.DecimalPrecision().DeclareModel()
pool.DecimalPrecision().Methods().CheckMainCurrencyRounding().DeclareMethod(
`CheckMainCurrencyRounding` ,
func (rs pool.DecimalPrecisionSet){
  //@api.constrains('digits')
  /*def _check_main_currency_rounding(self):
        if any(precision.name == 'Account' and
                tools.float_compare(self.env.user.company_id.currency_id.rounding, 10 ** - precision.digits, precision_digits=6) == -1
                for precision in self):
            raise ValidationError(_("You cannot define the decimal precision of 'Account' as greater than the rounding factor of the company's main currency"))
        return True
*/})
 
 }