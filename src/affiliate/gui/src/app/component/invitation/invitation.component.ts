import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { ApiService } from '../../service/api.service';
import { Ng4LoadingSpinnerService } from 'ng4-loading-spinner';

declare let PullToRefresh: any;

@Component({
  selector: 'app-invitation',
  templateUrl: './invitation.component.html',
  styleUrls: ['./invitation.component.css']
})
export class InvitationComponent implements OnInit {
  @ViewChild('invatationList') invatationList: ElementRef;
  ngAfterViewInit() {
    this.pullToRefreshRef = PullToRefresh.init({
      mainElement: this.invatationList.nativeElement,
      onRefresh: () => {
        return this.viewInvitation(this.params);
      }
    });
  }
  private pullToRefreshRef: any;
  private subscribeRef = null;
  private loaded: boolean = true;
  private params: any;
  invitationList: any[];
  constructor(
    private apiService: ApiService,
    private router: Router,
    private activeRoute: ActivatedRoute,
    private spinnerService: Ng4LoadingSpinnerService
  ) { }

  ngOnInit() {
    this.activeRoute.params.subscribe(params => {
      //console.log(params);
      if (!params.address) {
        console.log("no address");
        return;
      }
      this.params = params;
      this.viewInvitation(params);
    });
  }
  ngOnDestroy() {
    this.spinnerService.hide();
    if (this.subscribeRef) {
      this.subscribeRef.unsubscribe();
    }
    if (this.pullToRefreshRef) {
      this.pullToRefreshRef.destroy();
    }
  }
  viewInvitation(params) {
    return new Promise((resolve, reject) => {
      if (!this.loaded) {
        this.spinnerService.show();
      }
      this.subscribeRef = this.apiService.post("/code/my-invitation/", params).subscribe(res => {
        console.log(res)
        this.invitationList = res.list;
        this.spinnerService.hide();
        this.loaded = true;
        resolve();
      }, err => {
        alert(err);
        this.spinnerService.hide();
        this.loaded = true;
        reject();
      })
    })

  }
}
