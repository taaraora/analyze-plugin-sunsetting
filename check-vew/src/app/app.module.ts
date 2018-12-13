import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { CheckViewComponentComponent } from './check-view-component/check-view-component.component';

@NgModule({
  declarations: [
    AppComponent,
    CheckViewComponentComponent
  ],
  imports: [
    BrowserModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
