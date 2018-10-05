## Version 1.5.1 (Release date: 2018-10-05)
([59a68dd](https://github.com/nytimes/video-transcoding-api/commit/59a68dd)) readme: move logo up 





## Version 1.5.0 (Release date: 2018-10-05)

([201804b](https://github.com/nytimes//commit/201804b)) deploy: app name isn't secret 

([2bcb17c](https://github.com/nytimes//commit/2bcb17c)) travis-script/deploy: update to use generic strat 

([da62d10](https://github.com/nytimes//commit/da62d10)) Add build script back 

([6cfcb50](https://github.com/nytimes//commit/6cfcb50)) travis: sudo is not required 

([aa5741a](https://github.com/nytimes//commit/aa5741a)) Remove drone stuff and setup travis to trigger deployment on drone 

([7c74a75](https://github.com/nytimes//commit/7c74a75)) Dockerfile: install ca-certificates 


([ede9d11](https://github.com/nytimes//commit/ede9d11)) removing tests now that manifest generation occurs with the encode, everything should be good 

([7e55b08](https://github.com/nytimes//commit/7e55b08)) fixing path issue for hls outputs 

([e9fe31b](https://github.com/nytimes//commit/e9fe31b)) manifest generation is now part of the encode step 


([6ba10f0](https://github.com/nytimes//commit/6ba10f0)) travis: fix build 

([779cd3a](https://github.com/nytimes//commit/779cd3a)) drone: fix build 

([f61b061](https://github.com/nytimes//commit/f61b061)) drone: set GO111MODULE to on 

([6f11490](https://github.com/nytimes//commit/6f11490)) Update build script 

([18107aa](https://github.com/nytimes//commit/18107aa)) Bye dep, welcome go mod 


([fb1392c](https://github.com/nytimes//commit/fb1392c)) remove expose instruction 

([be8ae34](https://github.com/nytimes//commit/be8ae34)) add dockerfile 

([d22ad46](https://github.com/nytimes//commit/d22ad46)) config: workaround for gofmt issue 

([d8b90aa](https://github.com/nytimes//commit/d8b90aa)) Makefile: add more params for golangci-lint 

([0a1ad60](https://github.com/nytimes//commit/0a1ad60)) config: fix gofmt 

([7251631](https://github.com/nytimes//commit/7251631)) Makefile: migrate from gometalinter to golangci-lint 

([51a30d6](https://github.com/nytimes//commit/51a30d6)) provider/bitmovin: remove unused constant 

([c6d8439](https://github.com/nytimes//commit/c6d8439)) drone: remove autoscaling alert config 

([2524e3b](https://github.com/nytimes//commit/2524e3b)) Makefile/gometalinter: disable gosec 

([73d17ea](https://github.com/nytimes//commit/73d17ea)) Update dependencies 






## Version 1.4.2 (Release date: 2018-04-19)
([69f9799](https://github.com/nytimes/video-transcoding-api/commit/69f9799)) provider/bitmovin: set StreamConditionsMode 

([f10fef2](https://github.com/nytimes/video-transcoding-api/commit/f10fef2)) Update bitmovin-go 


([fc43f30](https://github.com/nytimes/video-transcoding-api/commit/fc43f30)) Makefile: remove HTTP_ACCESS_LOG 

([c05302d](https://github.com/nytimes/video-transcoding-api/commit/c05302d)) bitmovin: encodes only existing audio for vp8 


([5b5fe4d](https://github.com/nytimes/video-transcoding-api/commit/5b5fe4d)) bitmovin:encodes audio only when a track is available 

([b1a44ac](https://github.com/nytimes/video-transcoding-api/commit/b1a44ac)) Update dependencies 

([cc0d223](https://github.com/nytimes/video-transcoding-api/commit/cc0d223)) Add support for specifying TwoPass on presets 

([b7fe798](https://github.com/nytimes/video-transcoding-api/commit/b7fe798)) Update dependencies 

([31b869e](https://github.com/nytimes/video-transcoding-api/commit/31b869e)) Remove swagger-ui stuff 

([79d2138](https://github.com/nytimes/video-transcoding-api/commit/79d2138)) db/types: some docs fixes 




## Version 1.4.1 (Release date: 2018-03-06)
([642021a](https://github.com/nytimes/video-transcoding-api/commit/642021a)) Update logging dependencies 

([57e3ff5](https://github.com/nytimes/video-transcoding-api/commit/57e3ff5)) Fix the build 

([0a14ef9](https://github.com/nytimes/video-transcoding-api/commit/0a14ef9)) Silence gizmo logger 

([7e47e62](https://github.com/nytimes/video-transcoding-api/commit/7e47e62)) swagger: make WithStatus ignore 0s 




## Version 1.4.0 (Release date: 2018-02-21)
([5382f42](https://github.com/nytimes/video-transcoding-api/commit/5382f42)) service: support sending access logs through logrus logger 

([4f2bd85](https://github.com/nytimes/video-transcoding-api/commit/4f2bd85)) Add SERVICE_NAME to log tags 

([cca255b](https://github.com/nytimes/video-transcoding-api/commit/cca255b)) travis: Go 1.10 




## Version 1.3.5 (Release date: 2018-02-13)
([63d0810](https://github.com/nytimes/video-transcoding-api/commit/63d0810)) vendor: update dependencies 

([e9c9751](https://github.com/nytimes/video-transcoding-api/commit/e9c9751)) Makefile: disable gotype on gometalinter 

([cce824c](https://github.com/nytimes/video-transcoding-api/commit/cce824c)) drone: fix downstream calls 




## Version 1.3.4 (Release date: 2018-01-30)
([3e2de79](https://github.com/nytimes/video-transcoding-api/commit/3e2de79)) Update dependencies 


([fad3c83](https://github.com/nytimes/video-transcoding-api/commit/fad3c83)) provider/bitmovin: add metadata info for webm and mov files 

([1d38677](https://github.com/nytimes/video-transcoding-api/commit/1d38677)) Gopkg: update bitmovin-go and add prune config 

([1df0ec5](https://github.com/nytimes/video-transcoding-api/commit/1df0ec5)) Makefile: use fast linters 


([97f11c8](https://github.com/nytimes/video-transcoding-api/commit/97f11c8)) provider/bitmovin: move creating encoding service out of a loop 

([eb843bc](https://github.com/nytimes/video-transcoding-api/commit/eb843bc)) provider/bitmovin: return error if any expected data is missing for job status 

([8676701](https://github.com/nytimes/video-transcoding-api/commit/8676701)) provider/bitmovin: add mp4 output files metadata to finished jobs 

([ebb8fd4](https://github.com/nytimes/video-transcoding-api/commit/ebb8fd4)) update vendor bitmovin sdk 




## Version 1.3.3 (Release date: 2018-01-09)
([68213de](https://github.com/nytimes/video-transcoding-api/commit/68213de)) Update dependencies 




## Version 1.3.2 (Release date: 2018-01-05)
([6f2219e](https://github.com/nytimes/video-transcoding-api/commit/6f2219e)) drone: fix deploy call 




## Version 1.3.1 (Release date: 2018-01-05)


([d3e03d7](https://github.com/nytimes/video-transcoding-api/commit/d3e03d7)) provider/bitmovin: add SourceInfo to finished jobs 

([d8d8845](https://github.com/nytimes/video-transcoding-api/commit/d8d8845)) Gopkg: use bitmovin-go@master 

([4e4179a](https://github.com/nytimes/video-transcoding-api/commit/4e4179a)) Gopkg: make go-redis spec more explicit 

([f3573c3](https://github.com/nytimes/video-transcoding-api/commit/f3573c3)) Update dependencies 

([1233a1d](https://github.com/nytimes/video-transcoding-api/commit/1233a1d)) provider/bitmovin: support full paths for mov output 


([0449ec6](https://github.com/nytimes/video-transcoding-api/commit/0449ec6)) adding mov support 

([fef0195](https://github.com/nytimes/video-transcoding-api/commit/fef0195)) provider/bitmovin: add output "folder" to job status 

([4edb190](https://github.com/nytimes/video-transcoding-api/commit/4edb190)) provider/bitmovin: place output files inside folders with job ID 


([bf03fb3](https://github.com/nytimes/video-transcoding-api/commit/bf03fb3)) drone: fix autoscaling group name in notification config 

([51ae790](https://github.com/nytimes/video-transcoding-api/commit/51ae790)) provider/bitmovin: refactor JobStatus and add some extra info 

([6cc290d](https://github.com/nytimes/video-transcoding-api/commit/6cc290d)) Update dependencies 

([91a3de1](https://github.com/nytimes/video-transcoding-api/commit/91a3de1)) provider/bitmovin: add progress 


([5be7f61](https://github.com/nytimes/video-transcoding-api/commit/5be7f61)) drone: change deployment config 

([070bbc4](https://github.com/nytimes/video-transcoding-api/commit/070bbc4)) provider/bitmovin: refactor s3 url parsing 

([2cbba3d](https://github.com/nytimes/video-transcoding-api/commit/2cbba3d)) provider/bitmovin: fix metalinter violation 


([908c910](https://github.com/nytimes/video-transcoding-api/commit/908c910)) dependency management go inflates the amount of lines of code i have changed 

([e12421b](https://github.com/nytimes/video-transcoding-api/commit/e12421b)) changes to our api client so new version 

([d26abab](https://github.com/nytimes/video-transcoding-api/commit/d26abab)) removing debug code 

([e3b60af](https://github.com/nytimes/video-transcoding-api/commit/e3b60af)) do not check in coverage files 

([0f55e37](https://github.com/nytimes/video-transcoding-api/commit/0f55e37)) more test coverage 

([d0f4237](https://github.com/nytimes/video-transcoding-api/commit/d0f4237)) adding vp8, fixing tests 

([fd0025e](https://github.com/nytimes/video-transcoding-api/commit/fd0025e)) some debug code 

([4de37e5](https://github.com/nytimes/video-transcoding-api/commit/4de37e5)) adding vp8 


([045cb4d](https://github.com/nytimes/video-transcoding-api/commit/045cb4d)) gofmt -s -w . 

([6511c74](https://github.com/nytimes/video-transcoding-api/commit/6511c74)) db/redis: update version of go-redis 

([7573469](https://github.com/nytimes/video-transcoding-api/commit/7573469)) travis: remove Go 1.8 

([02437e4](https://github.com/nytimes/video-transcoding-api/commit/02437e4)) Another shot at vendoring 

([80c9401](https://github.com/nytimes/video-transcoding-api/commit/80c9401)) Fix gops agent 

([0645448](https://github.com/nytimes/video-transcoding-api/commit/0645448)) Makefile: remove unused 

([d915e08](https://github.com/nytimes/video-transcoding-api/commit/d915e08)) Revert "derp" 

([bf8826d](https://github.com/nytimes/video-transcoding-api/commit/bf8826d)) Revert "Makefile: no need for go get anymore" 

([f4dc40b](https://github.com/nytimes/video-transcoding-api/commit/f4dc40b)) Revert "bin/build: no need for go get" 

([101e70a](https://github.com/nytimes/video-transcoding-api/commit/101e70a)) Revert "bin/build: properly support vendoring" 

([39e8e74](https://github.com/nytimes/video-transcoding-api/commit/39e8e74)) bin/build: properly support vendoring 

([f1e77f8](https://github.com/nytimes/video-transcoding-api/commit/f1e77f8)) bin/build: no need for go get 

([7eaae8b](https://github.com/nytimes/video-transcoding-api/commit/7eaae8b)) Makefile: no need for go get anymore 

([33c43b5](https://github.com/nytimes/video-transcoding-api/commit/33c43b5)) derp 

([f893fb3](https://github.com/nytimes/video-transcoding-api/commit/f893fb3)) Revert "travis: add hack for Go 1.9" 

([cc27f73](https://github.com/nytimes/video-transcoding-api/commit/cc27f73)) travis: add hack for Go 1.9 

([5632c22](https://github.com/nytimes/video-transcoding-api/commit/5632c22)) travis: Go 1.9 

([ee7ff0c](https://github.com/nytimes/video-transcoding-api/commit/ee7ff0c)) Remove go 1.9 from travis 

([c51998f](https://github.com/nytimes/video-transcoding-api/commit/c51998f)) Makefile: set GOROOT before invoking swagger generate 

([a2b5a70](https://github.com/nytimes/video-transcoding-api/commit/a2b5a70)) travis: run tests on Go1.9rc1 

([60c4590](https://github.com/nytimes/video-transcoding-api/commit/60c4590)) Use official sdhook repo 

([f696ba7](https://github.com/nytimes/video-transcoding-api/commit/f696ba7)) Fix logrus import path 

([93b2edc](https://github.com/nytimes/video-transcoding-api/commit/93b2edc)) Makefile: use go build -i instead to install dependencies for lint 





## Version 1.3.0 (Release date: 2017-06-13)

([7695f83](https://github.com/nytimes/video-transcoding-api/commit/7695f83)) Add default Hybrik PresetPath and ComplianceDate settings to config 

([0641150](https://github.com/nytimes/video-transcoding-api/commit/0641150)) Allow Hybrik preset path to be optionaly configurable via env config var 

([4532d91](https://github.com/nytimes/video-transcoding-api/commit/4532d91)) Fix code formatting issues 

([14b8876](https://github.com/nytimes/video-transcoding-api/commit/14b8876)) Update license information 

([86b8f30](https://github.com/nytimes/video-transcoding-api/commit/86b8f30)) db/redis: truncate time to milliseconds 

([b3e136a](https://github.com/nytimes/video-transcoding-api/commit/b3e136a)) Hybrik provider: Only apply StreamingParams.SegmentDuration to HLS outputs 

([a4eafc7](https://github.com/nytimes/video-transcoding-api/commit/a4eafc7)) Fix gometalinter warnings 


([d126217](https://github.com/nytimes/video-transcoding-api/commit/d126217)) Parallelize preset retrieval on job creation 

([8cc97c6](https://github.com/nytimes/video-transcoding-api/commit/8cc97c6)) Switch from encoding-wrapper/hybrik to hybrik-sdk-go package 

([286b126](https://github.com/nytimes/video-transcoding-api/commit/286b126)) Add Hybrik to README 

([d354cbc](https://github.com/nytimes/video-transcoding-api/commit/d354cbc)) Add support for Hybrik as a transcoding provider 

([c49c678](https://github.com/nytimes/video-transcoding-api/commit/c49c678)) fix Readme 


([06e5856](https://github.com/nytimes/video-transcoding-api/commit/06e5856)) forgot to add Destination field in README 

([3896683](https://github.com/nytimes/video-transcoding-api/commit/3896683)) adding more tests 

([4d9819e](https://github.com/nytimes/video-transcoding-api/commit/4d9819e)) removing debug code 

([4ae3275](https://github.com/nytimes/video-transcoding-api/commit/4ae3275)) added ability to set output s3 destination, decoupling from input.  about to add more input types 

([145bf16](https://github.com/nytimes/video-transcoding-api/commit/145bf16)) adding debugging code will remove later 

([dd86b3c](https://github.com/nytimes/video-transcoding-api/commit/dd86b3c)) adding ability to specify encoding version 





## Version 1.2.1 (Release date: 2017-04-04)
([1212d45](https://github.com/nytimes/video-transcoding-api/commit/1212d45)) provider/encodingcom: fix bug in newly introduced adjustSize function 




## Version 1.2.0 (Release date: 2017-04-04)
([9cbe954](https://github.com/nytimes/video-transcoding-api/commit/9cbe954)) Fix gometalinter violations 

([5776807](https://github.com/nytimes/video-transcoding-api/commit/5776807)) provider/encodingcom: provide the proper information on JobStatus 


([41ae17b](https://github.com/nytimes/video-transcoding-api/commit/41ae17b)) new fields for README 

([05dc0ea](https://github.com/nytimes/video-transcoding-api/commit/05dc0ea)) adding flexibility for s3 storage region and the ability to encode in any region of aws or gcp 


([6f5cf95](https://github.com/nytimes/video-transcoding-api/commit/6f5cf95)) travis: use .x notation for Go 1.8 


([8eaa73e](https://github.com/nytimes/video-transcoding-api/commit/8eaa73e)) code cleanup 

([6deee07](https://github.com/nytimes/video-transcoding-api/commit/6deee07)) removing debugging code 

([1d8e517](https://github.com/nytimes/video-transcoding-api/commit/1d8e517)) linter caught some dead code and other things, fixing this 

([e93afc3](https://github.com/nytimes/video-transcoding-api/commit/e93afc3)) code cleanup 

([b12c739](https://github.com/nytimes/video-transcoding-api/commit/b12c739)) more test coverage and code cleanup 

([0c1c815](https://github.com/nytimes/video-transcoding-api/commit/0c1c815)) adding test code 

([da1a24e](https://github.com/nytimes/video-transcoding-api/commit/da1a24e)) travis: update go 1.8 to rc3 

([1b3bd2d](https://github.com/nytimes/video-transcoding-api/commit/1b3bd2d)) adding Bitmovin to readme 

([ccb4321](https://github.com/nytimes/video-transcoding-api/commit/ccb4321)) jobs submit but does not respond with proper json, debugging code in here for now, manifest gen works for hls 

([4af9b33](https://github.com/nytimes/video-transcoding-api/commit/4af9b33)) changing import directory 

([7a0a658](https://github.com/nytimes/video-transcoding-api/commit/7a0a658)) need to register bitmovin 

([2900b79](https://github.com/nytimes/video-transcoding-api/commit/2900b79)) fleshing out rest of methods 

([ec748c5](https://github.com/nytimes/video-transcoding-api/commit/ec748c5)) interstitial commit 

([29272e9](https://github.com/nytimes/video-transcoding-api/commit/29272e9)) change in logic, adding a lot of steps 

([e2c10d9](https://github.com/nytimes/video-transcoding-api/commit/e2c10d9)) adding environment variables for bitmovin aws credentials 

([2fbc577](https://github.com/nytimes/video-transcoding-api/commit/2fbc577)) s3 url parsing 

([3a05910](https://github.com/nytimes/video-transcoding-api/commit/3a05910)) interstitial commit 

([50b89f3](https://github.com/nytimes/video-transcoding-api/commit/50b89f3)) interstitial commit 

([595dce8](https://github.com/nytimes/video-transcoding-api/commit/595dce8)) implementing the interface 

([a34faef](https://github.com/nytimes/video-transcoding-api/commit/a34faef)) adding rest of config for bitmovin clinet 

([02e97d3](https://github.com/nytimes/video-transcoding-api/commit/02e97d3)) initial config changes for bitmovin 

([003fcd1](https://github.com/nytimes/video-transcoding-api/commit/003fcd1)) adding stubs for functions, capapbilities and first test 

([6f424e6](https://github.com/nytimes/video-transcoding-api/commit/6f424e6)) initial commit 





## Version 1.1.2 (Release date: 2017-01-25)

([f32e347](https://github.com/nytimes/video-transcoding-api/commit/f32e347)) Changes the Zencoder wrapper to use the job status found in the jobDetails for the JobState 

([58b5e83](https://github.com/nytimes/video-transcoding-api/commit/58b5e83)) travis: update go 1.8 to rc2 

([9357eaa](https://github.com/nytimes/video-transcoding-api/commit/9357eaa)) db/redis: update go-redis 

([17730c0](https://github.com/nytimes/video-transcoding-api/commit/17730c0)) drone: run integration tests after deploying 

([fa7488e](https://github.com/nytimes/video-transcoding-api/commit/fa7488e)) travis: update Go 1.8 to rc1 

([ba8e873](https://github.com/nytimes/video-transcoding-api/commit/ba8e873)) readme: don't use an actual IP address in example 

([43e243e](https://github.com/nytimes/video-transcoding-api/commit/43e243e)) Update swagger.json 

([f4bd0eb](https://github.com/nytimes/video-transcoding-api/commit/f4bd0eb)) doc: include zencoder 




## Version 1.1.1 (Release date: 2017-01-06)

([98ba09f](https://github.com/nytimes/video-transcoding-api/commit/98ba09f)) travis: use .x syntax to ensure latest 1.7 

([24d9528](https://github.com/nytimes/video-transcoding-api/commit/24d9528)) db/redis/storage: support float64 




## Version 1.1.0 (Release date: 2016-12-22)

([ec42742](https://github.com/nytimes/video-transcoding-api/commit/ec42742)) encodingcom: Returns converted file size 

([12b0a6d](https://github.com/nytimes/video-transcoding-api/commit/12b0a6d)) Update gops 

([0c1152c](https://github.com/nytimes/video-transcoding-api/commit/0c1152c)) encodingcom: Returns converted file size 





## Version 1.0.9 (Release date: 2016-12-19)



## Version 1.0.8-rc (Release date: 2016-12-19)

([071a2d8](https://github.com/nytimes/video-transcoding-api/commit/071a2d8)) preset: avoid creating PresetMap when ProviderMapping is empty 

([314930b](https://github.com/nytimes/video-transcoding-api/commit/314930b)) preset: improve if statement 




## Version 1.0.7 (Release date: 2016-12-16)



## Version 1.0.6-rc (Release date: 2016-12-16)

([d863498](https://github.com/nytimes//commit/d863498)) Add filesize rendition info to Zencoder 


([6de4a74](https://github.com/nytimes//commit/6de4a74)) preset: bubble up the error when creating a preset 

([39c0992](https://github.com/nytimes//commit/39c0992)) service/preset: remove logging and fix comments 

([90d24dc](https://github.com/nytimes//commit/90d24dc)) service/presetmap: create or update existent presetmap when creating new presets 




## Version 1.0.5 (Release date: 2016-12-08)
([c53ea3d](https://github.com/nytimes/video-transcoding-api/commit/c53ea3d)) provider/elementalconductor: fix compatibility with encoding-wrapper 

([6e951cc](https://github.com/nytimes/video-transcoding-api/commit/6e951cc)) travis: add Go 1.8beta1 

([0295924](https://github.com/nytimes/video-transcoding-api/commit/0295924)) swagger: add test to increase package coverage 


([c60a8f3](https://github.com/nytimes/video-transcoding-api/commit/c60a8f3)) zencoder: use constants for Job State 

([44e3bd5](https://github.com/nytimes/video-transcoding-api/commit/44e3bd5)) zencoder: set progress to 100 when job status is finished (fixes #170) 




## Version 1.0.4 (Release date: 2016-12-05)

([b220d80](https://github.com/nytimes/video-transcoding-api/commit/b220d80)) zencoder: bugfix on duration being reported 


([4566b29](https://github.com/nytimes/video-transcoding-api/commit/4566b29)) Revert "encodingcom: remove dead code when creating a encoding.com Format" 




## Version 1.0.3 (Release date: 2016-12-02)

([db75994](https://github.com/nytimes/video-transcoding-api/commit/db75994)) Protect against possible stray colon in error message 

([a56d8da](https://github.com/nytimes/video-transcoding-api/commit/a56d8da)) Revert "Populate status message" 


([4df203f](https://github.com/nytimes/video-transcoding-api/commit/4df203f)) Populate job status with detailed status message 




## Version 1.0.2-rc (Release date: 2016-12-01)
([98d5b5a](https://github.com/nytimes/video-transcoding-api/commit/98d5b5a)) Makefile: use CI_TAG in `make live` when available 

([d0c5e27](https://github.com/nytimes/video-transcoding-api/commit/d0c5e27)) travis: update Go 




## Version 1.0.1-rc (Release date: 2016-12-01)

([a6112be](https://github.com/nytimes/video-transcoding-api/commit/a6112be)) zencoder: consider finished outputs with no format and m3u8 suffix as m3u8 container (refs #161) 


([d68b5e2](https://github.com/nytimes/video-transcoding-api/commit/d68b5e2)) provider/zencoder: use GetVodUsage in Healthcheck 




([58bd6e8](https://github.com/nytimes/video-transcoding-api/commit/58bd6e8)) makefile: fix variable replacement for makefile scheme 

([766e434](https://github.com/nytimes/video-transcoding-api/commit/766e434)) build: detach stg and prod deployment by using 'rc' on tag name 

([b3d1b63](https://github.com/nytimes/video-transcoding-api/commit/b3d1b63)) encodingcom: remove dead code when creating a encoding.com Format 




## Version 1.0.0 (Release date: 2016-11-23)



## Version 0.1.6 (Release date: 2016-11-23)
([dbf6d29](https://github.com/nytimes/video-transcoding-api/commit/dbf6d29)) zencoder: fix hls path (close #157) 




## Version 0.1.5 (Release date: 2016-11-23)
([bf7521b](https://github.com/nytimes/video-transcoding-api/commit/bf7521b)) zencoder: avoid concatenating 'hls' path to hls output (refs #157) 





## Version 0.1.4 (Release date: 2016-11-23)




## Version 0.1.2 (Release date: 2016-11-21)

([f3e2435](https://github.com/nytimes/video-transcoding-api/commit/f3e2435)) zencoder: add PrepareForSegmenting: 'hls' for mp4's that matches with HLS 

([e967ab0](https://github.com/nytimes/video-transcoding-api/commit/e967ab0)) zencoder: fix golint complain 

([2a0a755](https://github.com/nytimes/video-transcoding-api/commit/2a0a755)) zencoder: bugfix on isOutputCompatible() method 

([2d69d00](https://github.com/nytimes/video-transcoding-api/commit/2d69d00)) zencoder: raise errors gracefully 

([59e6e75](https://github.com/nytimes/video-transcoding-api/commit/59e6e75)) zencoder: reuse mp4 outputs for transmuxing hls outputs (close #151) 


([7f0d41c](https://github.com/nytimes/video-transcoding-api/commit/7f0d41c)) Send logging and error reporting via agent 




## Version 0.1.1 (Release date: 2016-11-16)

([225d9ff](https://github.com/nytimes/video-transcoding-api/commit/225d9ff)) zencoder: make all zencoder uploads public 


([37acd1c](https://github.com/nytimes/video-transcoding-api/commit/37acd1c)) zencoder: normalize hls output names based on encoding.com implementation 


([8666837](https://github.com/nytimes/video-transcoding-api/commit/8666837)) db: remove unused struct from stub_test 

([d6a1595](https://github.com/nytimes/video-transcoding-api/commit/d6a1595)) db/redis: add structs to stub_test to avoid db dependency on redis storage 

([a70a602](https://github.com/nytimes/video-transcoding-api/commit/a70a602)) db/redis: add more tests for FieldMap() method 

([4231353](https://github.com/nytimes/video-transcoding-api/commit/4231353)) db/redis: add test for FieldMap() method 