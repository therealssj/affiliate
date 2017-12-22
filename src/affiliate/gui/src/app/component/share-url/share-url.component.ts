import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { ApiService } from '../../service/api.service';
import { Ng4LoadingSpinnerService } from 'ng4-loading-spinner';

@Component({
  selector: 'app-share-url',
  templateUrl: './share-url.component.html',
  styleUrls: ['./share-url.component.css']
})
export class ShareUrlComponent implements OnInit {
  private subscribeRef = null;
  private buyUrl = "";
  private joinUrl = "";
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
        console.log("no address");
        return;
      }
      this.generateInfo(params);
    }); 
  }
  ngOnDestroy() {
    this.spinnerService.hide();
    if(this.subscribeRef) {
      this.subscribeRef.unsubscribe();
    }
  }
  generateInfo(params: any) {
    this.spinnerService.show();
    this.subscribeRef = this.apiService.post("/code/generate/", params).subscribe(res => {
      //console.log(res)
      this.buyUrl = res.buyUrl;
      this.joinUrl = res.joinUrl;
    }, err => {
      alert(err);
      this.spinnerService.hide();
    }, ()=>{
      this.spinnerService.hide();
    })
  }

}
