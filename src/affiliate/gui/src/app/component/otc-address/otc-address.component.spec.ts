import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OtcAddressComponent } from './otc-address.component';

describe('OtcAddressComponent', () => {
  let component: OtcAddressComponent;
  let fixture: ComponentFixture<OtcAddressComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OtcAddressComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OtcAddressComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
