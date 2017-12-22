import { RouterModule } from '@angular/router';

import { AppComponent } from './app.component';
import { IndexComponent } from './component/index/index.component';
import { ShareUrlComponent } from './component/share-url/share-url.component';
import { InvitationComponent } from './component/invitation/invitation.component';

export const appRoutes=[
	{
		path:'',
		component:IndexComponent
	},
	{
		path:"shareUrl",
		component:ShareUrlComponent
	},
	{
		path:"invitation",
		component:InvitationComponent
	},
	{
		path:'**',//fallback router must in the last
		component:IndexComponent
	}
];
