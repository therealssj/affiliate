import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OtcStatusComponent } from './otc-status.component';

describe('OtcStatusComponent', () => {
  let component: OtcStatusComponent;
  let fixture: ComponentFixture<OtcStatusComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OtcStatusComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OtcStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
