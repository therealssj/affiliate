import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations'; 
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms'; //引入表单模块
//ui
import {
  MatButtonModule, 
  MatCheckboxModule, 
  MatInputModule
} from '@angular/material';
import {MatListModule} from '@angular/material/list'
//loading
import { Ng4LoadingSpinnerModule } from 'ng4-loading-spinner';

import { AppComponent } from './app.component';

import { RouterModule } from '@angular/router';
import { appRoutes } from './app.routes';
import { CodeComponent } from './component/code/code.component';
import { ShareUrlComponent } from './component/share-url/share-url.component';
import { InvitationComponent } from './component/invitation/invitation.component';

import { HttpModule } from '@angular/http';
import { ApiService } from "./service/api.service";

@NgModule({
  declarations: [
    AppComponent,
    CodeComponent,
    ShareUrlComponent,
    InvitationComponent
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
    RouterModule.forRoot(appRoutes, {useHash: true}),
    Ng4LoadingSpinnerModule.forRoot()
  ],
  providers: [ApiService],
  exports: [MatButtonModule, MatCheckboxModule, MatInputModule, MatListModule],
  bootstrap: [AppComponent]
})
export class AppModule { }
