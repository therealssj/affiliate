import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/timeout';

declare let config:any;

@Injectable()
export class ApiService {

  private domain = config.domain; 

  constructor(private http: Http) { }

  get(url, options = null) {
    return this.http.get(this.getUrl(url, options)).timeout(15000)
      .map((res: any) => {
        let obj = res.json();
        return obj.code === 0 ? obj.data : this.throwError(obj.errmsg);
      })
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }
  
  post(url, options = {}) {
    return this.http.post(this.getUrl(url), this.getQueryString(options), this.returnRequestOptions()).timeout(15000)
      .map((res: any) => {
        let obj = res.json();
        return obj.code === 0 ? obj.data : this.throwError(obj.errmsg);
      })
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  throwError(e:any) {
    throw e;
  }

  private getHeaders() {
    const headers = new Headers();
    headers.append('Content-Type', 'application/x-www-form-urlencoded');
    return headers;
  }

  returnRequestOptions() {
    const options = new RequestOptions();

    options.headers = this.getHeaders();

    return options;
  }

  private getQueryString(parameters = null) {
    if (!parameters) {
      return '';
    }

    return Object.keys(parameters).reduce((array,key) => {
      array.push(key + '=' + encodeURIComponent(parameters[key]));
      return array;
    }, []).join('&');
  }

  private getUrl(url, options = null) {
    if(options == null){
      return this.domain + url;
    }
    return this.domain + url + '?' + this.getQueryString(options);
  }
}
