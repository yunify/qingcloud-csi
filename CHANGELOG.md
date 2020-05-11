<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [v1.2.0](#v120)
  - [Changelog since v1.1.0](#changelog-since-v110)
    - [Features](#features)
    - [Bug Fixed](#bug-fixed)
- [v1.1.1](#v111)
  - [Changelog since v1.1.0](#changelog-since-v110-1)
    - [Features](#features-1)
- [v1.1.0](#v110)
  - [Changelog since v0.2.1](#changelog-since-v021)
    - [Features](#features-2)
    - [Bug Fixed](#bug-fixed-1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

<!-- NEW RELEASE NOTES ENTRY -->

# v1.2.0
## Changelog since v1.1.0
### Features
* Update CSI v1.2.0 and QingCloud SDK ([#164](https://github.com/yunify/qingcloud-csi/pull/164), [@wnxn](https://github.com/wnxn))
* Add retry limiter when detach volume ([#160](https://github.com/yunify/qingcloud-csi/pull/160), [@wnxn](https://github.com/wnxn))
* Support deploy on Kubernetes v1.16 ([#168](https://github.com/yunify/qingcloud-csi/pull/168), [@wnxn](https://github.com/wnxn))

### Bug Fixed
* Update snapshot content source ([#164](https://github.com/yunify/qingcloud-csi/pull/164), [@wnxn](https://github.com/wnxn))

# v1.1.1
## Changelog since v1.1.0
### Features
* Support more volume types ([#170](https://github.com/yunify/qingcloud-csi/pull/170), [@wnxn](https://github.com/wnxn))

# v1.1.0
## Changelog since v0.2.1
### Features
* Update to CSI v1.1.0 ([#62](https://github.com/yunify/qingcloud-csi/pull/62), [@wnxn](https://github.com/wnxn))
* Add snapshot creation and deletion ([#72](https://github.com/yunify/qingcloud-csi/pull/72), [@wnxn](https://github.com/wnxn))
* Add restoring volume from snapshot ([#77](https://github.com/yunify/qingcloud-csi/pull/77), [@wnxn](https://github.com/wnxn))
* Add offline volume expansion ([#85](https://github.com/yunify/qingcloud-csi/pull/85), [@wnxn](https://github.com/wnxn))
* Add to retrieve volume stats ([#88](https://github.com/yunify/qingcloud-csi/pull/88), [@wnxn](https://github.com/wnxn))
* Add stage build image ([#89](https://github.com/yunify/qingcloud-csi/pull/89), [@wnxn](https://github.com/wnxn))
* Add topology awareness ([#102](https://github.com/yunify/qingcloud-csi/pull/102), [@wnxn](https://github.com/wnxn))
* Replace Dep with Go mod ([#103](https://github.com/yunify/qingcloud-csi/pull/103), [@wnxn](https://github.com/wnxn))
* Add to bind tags and created resources ([#106](https://github.com/yunify/qingcloud-csi/pull/106), [@wnxn](https://github.com/wnxn))
* Use Dockerhub and remove QingCloud Dockerhub secret ([#107](https://github.com/yunify/qingcloud-csi/pull/107), [@wnxn](https://github.com/wnxn))
* Add mutex to handle concurrency requests ([#109](https://github.com/yunify/qingcloud-csi/pull/109), [@wnxn](https://github.com/wnxn))
* Add volume cloning ([#111](https://github.com/yunify/qingcloud-csi/pull/111), [@wnxn](https://github.com/wnxn))
* Add volume name prefix flag ([#113](https://github.com/yunify/qingcloud-csi/pull/113), [@wnxn](https://github.com/wnxn))
* Add maximum retry interval flag ([#116](https://github.com/yunify/qingcloud-csi/pull/116), [@wnxn](https://github.com/wnxn))
* Add Guarantee class of QoS ([#123](https://github.com/yunify/qingcloud-csi/pull/123), [@wnxn](https://github.com/wnxn))
* Update volume attachment map ([#137](https://github.com/yunify/qingcloud-csi/pull/137), [@wnxn](https://github.com/wnxn))
* Support more volume and instance types ([#139](https://github.com/yunify/qingcloud-csi/pull/139), [@wnxn](https://github.com/wnxn))

### Bug Fixed
* Fix cannot find device path after volume attached ([#133](https://github.com/yunify/qingcloud-csi/pull/133), [@wnxn](https://github.com/wnxn))