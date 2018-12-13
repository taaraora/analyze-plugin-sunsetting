import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CheckViewComponentComponent } from './check-view-component.component';

describe('CheckViewComponentComponent', () => {
  let component: CheckViewComponentComponent;
  let fixture: ComponentFixture<CheckViewComponentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CheckViewComponentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CheckViewComponentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
