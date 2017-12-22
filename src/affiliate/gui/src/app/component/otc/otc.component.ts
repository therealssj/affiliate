import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../service/api.service';
import { Router } from '@angular/router';
import { FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-otc',
  templateUrl: './otc.component.html',
  styleUrls: ['./otc.component.css']
})
export class OtcComponent implements OnInit {

  private modal = {
    address: "",
    currencyType: "BTC"
  }

  constructor(
    private apiService: ApiService,
    private router: Router
  ) { }

  ngOnInit() { }

  onGet() {
    console.log(this.modal)
    this.router.navigate(['/otcAddress', { address: this.modal.address, currencyType: this.modal.currencyType, ref: "" }]);
  }
  onCheck(){
    this.router.navigate(['/otcStatus', { address: this.modal.address }]);
  }

}
