import { BrowserModule } from '@angular/platform-browser';
import {Injector, NgModule} from '@angular/core';

import { CheckViewComponentComponent } from './check-view-component/check-view-component.component';
import {createCustomElement} from "@angular/elements";

@NgModule({
  declarations: [
    CheckViewComponentComponent
  ],
  imports: [
    BrowserModule
  ],
  providers: [],
  entryComponents: [
    CheckViewComponentComponent
  ]
})
export class AppModule {
  constructor(private injector: Injector) {
  }

  ngDoBootstrap() {
    const nspCheckVewCE = createCustomElement(CheckViewComponentComponent, { injector: this.injector });
    customElements.define('node-sunsetting-plugin-check-view', nspCheckVewCE);
  }
}
