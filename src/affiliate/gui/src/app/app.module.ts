import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
//ui
import {
  MatButtonModule,
  MatCheckboxModule,
  MatInputModule,
  MatRadioModule
} from '@angular/material';
import { MatListModule } from '@angular/material/list';

//loading
import { Ng4LoadingSpinnerModule } from 'ng4-loading-spinner';

import { AppComponent } from './app.component';

import { RouterModule } from '@angular/router';
import { appRoutes } from './app.routes';
import { CodeComponent } from './component/code/code.component';
import { ShareUrlComponent } from './component/share-url/share-url.component';
import { InvitationComponent } from './component/invitation/invitation.component';

import { HttpModule } from '@angular/http';
import { ApiService } from './service/api.service';
import { OtcComponent } from './component/otc/otc.component';
import { OtcAddressComponent } from './component/otc-address/otc-address.component';
import { OtcStatusComponent } from './component/otc-status/otc-status.component';

@NgModule({
  declarations: [
    AppComponent,
    CodeComponent,
    ShareUrlComponent,
    InvitationComponent,
    OtcComponent,
    OtcAddressComponent,
    OtcStatusComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpModule,
    MatButtonModule,
    MatCheckboxModule,
    MatInputModule,
    MatListModule,
    MatRadioModule,
    RouterModule.forRoot(appRoutes, { useHash: false }),
    Ng4LoadingSpinnerModule.forRoot()
  ],
  providers: [ApiService],
  exports: [MatButtonModule, MatCheckboxModule, MatInputModule, MatRadioModule, MatListModule],
  bootstrap: [AppComponent]
})
export class AppModule { }
