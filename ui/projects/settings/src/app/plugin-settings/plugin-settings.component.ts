  import { Component, ViewEncapsulation, Input, ElementRef } from '@angular/core';
  import { HttpClient } from '@angular/common/http';
  import { FormBuilder, FormGroup, Validators } from '@angular/forms';

  import { Config } from '../models/models'

@Component({
  // selector: 'app-plugin-settings',
  templateUrl: './plugin-settings.component.html',
  styleUrls: ['./plugin-settings.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class PluginSettingsComponent {

  @Input('pluginConfig')
  set setConfig(pluginConfig: string) {
    try {
      let config: Config = JSON.parse(pluginConfig);
      this.defaultPluginConfig = config;
      this.userPluginConfig = this.formBuilder.group({
        executionInterval: [config.executionInterval, [ Validators.min(60), Validators.max(86400) ]]
      });
      this.diffConfigs();
    } catch (err) {
      console.error(err);
    }
  }

  private defaultPluginConfig: Config;
  public userPluginConfig: FormGroup;
  public enableSave: boolean;

  constructor (
    private http: HttpClient,
    private formBuilder: FormBuilder,
    private el: ElementRef
  ) { }


  public diffConfigs() {
    const enabled = (
      !(this.defaultPluginConfig.executionInterval == this.userPluginConfig.value.executionInterval) &&
      this.userPluginConfig.controls.executionInterval.valid
    )
    this.enableSave = enabled;
  }

  public reset() {
    this.userPluginConfig.controls.executionInterval.setValue(this.defaultPluginConfig.executionInterval);
    this.diffConfigs();
  }

  public save() {
    this.defaultPluginConfig.executionInterval = this.userPluginConfig.value.executionInterval;
    this.el.nativeElement.dispatchEvent(new CustomEvent('ConfigUpdate', {
      detail: this.defaultPluginConfig,
      bubbles: true
    }))
  }

  get executionInterval () {
    return this.userPluginConfig.get('executionInterval');
  }
}
