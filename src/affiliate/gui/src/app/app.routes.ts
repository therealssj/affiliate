import { RouterModule } from '@angular/router';

import { AppComponent } from './app.component';
import { CodeComponent } from './component/code/code.component';
import { ShareUrlComponent } from './component/share-url/share-url.component';
import { InvitationComponent } from './component/invitation/invitation.component';
import { OtcComponent } from './component/otc/otc.component';
import { OtcAddressComponent } from './component/otc-address/otc-address.component';
import { OtcStatusComponent } from './component/otc-status/otc-status.component';

export const appRoutes=[
	{
		path: '',
		redirectTo: "otc",
		pathMatch: 'full'
	},
	{
		path: "otc",
		component: OtcComponent
	},
	{
		path: "otcAddress",
		component: OtcAddressComponent
	},
	{
		path: "otcStatus",
		component: OtcStatusComponent
	},
	{
		path: "code",
		component: CodeComponent
	},
	{
		path: "shareUrl",
		component: ShareUrlComponent
	},
	{
		path: "invitation",
		component: InvitationComponent
	},
	{
		path: '**',//fallback router must in the last
		redirectTo: "otc",
		pathMatch: 'full'
	}
];
