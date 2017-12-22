import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OtcComponent } from './otc.component';

describe('OtcComponent', () => {
  let component: OtcComponent;
  let fixture: ComponentFixture<OtcComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OtcComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OtcComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
