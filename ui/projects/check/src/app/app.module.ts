import { BrowserModule } from '@angular/platform-browser';
import {CUSTOM_ELEMENTS_SCHEMA, NgModule} from '@angular/core';
import { CheckResultComponent } from './check-result/check-result.component';

@NgModule({
  declarations: [

  CheckResultComponent],
  imports: [
    BrowserModule
  ],
  providers: [],
  bootstrap: [],
  schemas: [
  CUSTOM_ELEMENTS_SCHEMA
]
})
export class AppModule { }
