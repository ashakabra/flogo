import { HttpModule } from "@angular/http";
import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { JiraQueryTicketActivityContribution} from "./activity";
import { WiServiceContribution} from "wi-studio/app/contrib/wi-contrib";

@NgModule({
    imports: [
      CommonModule,
      HttpModule
    ],
    providers: [
      {
         provide: WiServiceContribution,
         useClass: JiraQueryTicketActivityContribution
       }
    ]
  })
  
  export default class JiraQueryTicketActivityModule {
  
  }