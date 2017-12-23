import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../service/api.service';
import { Router } from '@angular/router';
import { FormControl, Validators } from '@angular/forms';
import { URLSearchParams } from '@angular/http';

declare let window;

@Component({
  selector: 'app-otc',
  templateUrl: './otc.component.html',
  styleUrls: ['./otc.component.css']
})
export class OtcComponent implements OnInit {

  modal = {
    address: "",
    currencyType: "BTC"
  }

  constructor(
    private apiService: ApiService,
    private router: Router
  ) { }

  ngOnInit() { }
  getUrlParams(name){
    //must be remove the "?", or the first key will start with '?'
    let searchParams = new URLSearchParams(window.location.search.replace(/^\?/, ""));
    return searchParams.get(name);
  }
  onGet() {
    console.log(this.modal)
    this.router.navigate(['/otcAddress', { address: this.modal.address, currencyType: this.modal.currencyType, ref: this.getUrlParams("ref") }]);
  }
  onCheck(){
    this.router.navigate(['/otcStatus', { address: this.modal.address }]);
  }

}
