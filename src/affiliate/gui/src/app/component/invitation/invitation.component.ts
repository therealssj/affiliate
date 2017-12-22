import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { ApiService } from '../../service/api.service';
import { Ng4LoadingSpinnerService } from 'ng4-loading-spinner';

@Component({
  selector: 'app-invitation',
  templateUrl: './invitation.component.html',
  styleUrls: ['./invitation.component.css']
})
export class InvitationComponent implements OnInit {
  private subscribeRef = null;
  private invitationList = [];
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
      this.viewInvitation(params);
    });    
  }
  ngOnDestroy() {
    this.spinnerService.hide();
    if(this.subscribeRef) {
      this.subscribeRef.unsubscribe();
    }
  }
  viewInvitation(params){
    this.spinnerService.show();
    this.subscribeRef = this.apiService.post("/code/my-invitation/", params).subscribe(res => {
      console.log(res)
      this.invitationList = res.list;
    }, err => {
      alert(err);
      this.spinnerService.hide();
    }, ()=>{
      this.spinnerService.hide();
    })
  }
}
