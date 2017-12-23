import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../service/api.service';
import { Router } from '@angular/router';
import { Ng4LoadingSpinnerService } from 'ng4-loading-spinner';
import { FormControl, Validators } from '@angular/forms';
import { URLSearchParams } from '@angular/http';

declare let window:any;

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
  getUrlParams(name){
    //must be remove the "?", or the first key will start with '?'
    let searchParams = new URLSearchParams(window.location.search.replace(/^\?/, ""));
    return searchParams.get(name);
  }
  onGenerate() {
    this.router.navigate(['/shareUrl', { address: this.modal.address, ref: this.getUrlParams("ref") }]);
  }
  onViewInvitation(){
    this.router.navigate(['/invitation', { address: this.modal.address }]);
  }
}
