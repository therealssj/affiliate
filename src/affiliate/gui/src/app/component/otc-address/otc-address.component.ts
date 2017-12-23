import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { ApiService } from '../../service/api.service';
import { Ng4LoadingSpinnerService } from 'ng4-loading-spinner';

@Component({
  selector: 'app-otc-address',
  templateUrl: './otc-address.component.html',
  styleUrls: ['./otc-address.component.css']
})
export class OtcAddressComponent implements OnInit {

  private subscribeRef = null;
  address:string;
  currencyType:string;
  depositAddr:string;
  constructor(
    private apiService: ApiService,
    private router: Router, 
    private activeRoute: ActivatedRoute,
    private spinnerService: Ng4LoadingSpinnerService
  ) { }

  ngOnInit() {
    this.activeRoute.params.subscribe(params => {
      //console.log(params);
      if(!params.address || !params.currencyType){
        console.log("params error");
        return;
      }
      this.showInfo(params);
    }); 
  }
  ngOnDestroy() {
    this.spinnerService.hide();
    if(this.subscribeRef) {
      this.subscribeRef.unsubscribe();
    }
  }
  showInfo(params: any) {
    this.spinnerService.show();
    this.subscribeRef = this.apiService.post("/get-address/", params).subscribe(res => {
      //console.log(res)
      this.address = params.address;
      this.currencyType = params.currencyType;
      this.depositAddr = res.depositAddr;
      this.spinnerService.hide();
    }, err => {
      alert(err);
      this.spinnerService.hide();
    })
  }
}
