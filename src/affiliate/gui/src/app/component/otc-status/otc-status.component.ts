import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { ApiService } from '../../service/api.service';
import { Ng4LoadingSpinnerService } from 'ng4-loading-spinner';

@Component({
  selector: 'app-otc-status',
  templateUrl: './otc-status.component.html',
  styleUrls: ['./otc-status.component.css']
})
export class OtcStatusComponent implements OnInit {

  private subscribeRef = null;
  private updated:any;
  private gotDeposit:any;
  private currencyType:any;
  private depositAmount:any;
  private sendCoin:any;
  private coinAmount:any;

  constructor(
    private apiService: ApiService,
    private router: Router, 
    private activeRoute: ActivatedRoute,
    private spinnerService: Ng4LoadingSpinnerService
  ) { }

  ngOnInit() {
    this.activeRoute.params.subscribe(params => {
      //console.log(params);
      if(!params.address){
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
    this.subscribeRef = this.apiService.post("/check-status/", params).subscribe(res => {
      this.updated = res.updated;
      this.gotDeposit = res.gotDeposit;
      this.currencyType = res.currencyType;
      this.depositAmount = res.depositAmount;
      this.sendCoin = res.sendCoin;
      this.coinAmount = res.coinAmount;
      this.spinnerService.hide();
    }, err => {
      alert(err);
      this.spinnerService.hide();
    })
  }
}
