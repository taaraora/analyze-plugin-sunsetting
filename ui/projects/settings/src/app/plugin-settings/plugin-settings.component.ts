import { Component, ViewEncapsulation, Input } from '@angular/core';

@Component({
  // selector: 'app-plugin-settings',
  templateUrl: './plugin-settings.component.html',
  styleUrls: ['./plugin-settings.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class PluginSettingsComponent {

  @Input('pluginConfig')
  set setConfig(pluginConfig: string) {
    this.pluginConfig = JSON.parse(pluginConfig);
  }

  pluginConfig: any;

  constructor() { }
}
