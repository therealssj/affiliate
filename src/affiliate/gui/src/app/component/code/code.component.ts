import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../service/api.service';
import { Router } from '@angular/router';
import { Ng4LoadingSpinnerService } from 'ng4-loading-spinner';
import { FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-code',
  templateUrl: './code.component.html',
  styleUrls: ['./code.component.css']
})
export class CodeComponent implements OnInit {

  private desc = "";

  private defDesc = "The default desc";

  private modal = {
    address: ""
  }

  constructor(
    private apiService: ApiService,
    private router: Router,
    private spinnerService: Ng4LoadingSpinnerService
  ) { }

  ngOnInit() {
    this.apiService.get("/code/notice/").subscribe(res => {
      this.desc = res.desc;
    }, err => {
      //alert(err);
      this.desc = this.defDesc;
    })
  }

  onGenerate() {
    // this.apiService.get("http://192.168.238.1/mydata1").subscribe(res=>{
    //   console.log(res);
    // }, err=>{
    //   console.log(err);
    // })

    this.router.navigate(['/shareUrl', { address: this.modal.address, ref: "" }]);
    // if(this.walletChecker.hasError('required'))
    //   return;
    // this.loading = true;
    // this.spinnerService.show();

  }
  onViewInvitation(){
    this.router.navigate(['/invitation', { address: this.modal.address }]);
  }
  // getWalletErrMsg() {
  //   return this.walletChecker.hasError('required') ? 'You must enter a value' : '';
  // }

}
