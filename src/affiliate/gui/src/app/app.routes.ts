import { RouterModule } from '@angular/router';

import { AppComponent } from './app.component';
import { CodeComponent } from './component/code/code.component';
import { ShareUrlComponent } from './component/share-url/share-url.component';
import { InvitationComponent } from './component/invitation/invitation.component';

export const appRoutes=[
	{
		path: '',
		redirectTo: "code",
		pathMatch: 'full'
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
		redirectTo: "code",
		pathMatch: 'full'
	}
];
