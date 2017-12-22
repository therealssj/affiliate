import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../service/api.service';
import { Router } from '@angular/router';
import { Ng4LoadingSpinnerService } from 'ng4-loading-spinner';
import {FormControl, Validators} from '@angular/forms';

@Component({
  selector: 'app-index',
  templateUrl: './index.component.html',
  styleUrls: ['./index.component.css']
})
export class IndexComponent implements OnInit {
/**
 * 1.ajax
 * 2.router跳转 接受参数
 * 3.loading
 * 
 */
  private walletChecker = new FormControl('', [Validators.required]);
  
  private loading = false;

  private modal = {
    walletAddress: ""
  }

  constructor(
    private apiService: ApiService, 
    private router: Router,
    private spinnerService: Ng4LoadingSpinnerService
  ) { }

  ngOnInit() {
  }

  onGenerate(){
    // this.apiService.get("http://192.168.238.1/mydata1").subscribe(res=>{
    //   console.log(res);
    // }, err=>{
    //   console.log(err);
    // })

    //this.router.navigate(['/invitation', { path: "1" }]);
    // if(this.walletChecker.hasError('required'))
    //   return;
    this.loading = true;
    console.log(1)
    this.spinnerService.show();

  }

  // getWalletErrMsg() {
  //   return this.walletChecker.hasError('required') ? 'You must enter a value' : '';
  // }

}
