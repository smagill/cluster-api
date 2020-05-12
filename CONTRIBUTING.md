# Contributing Guidelines
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Contributor License Agreements](#contributor-license-agreements)
- [Finding Things That Need Help](#finding-things-that-need-help)
- [Contributing a Patch](#contributing-a-patch)
- [Backporting a Patch](#backporting-a-patch)
  - [Merge Approval](#merge-approval)
  - [Google Doc Viewing Permissions](#google-doc-viewing-permissions)
  - [Issue and Pull Request Management](#issue-and-pull-request-management)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

Read the following guide if you're interested in contributing to cluster-api.

## Contributor License Agreements

We'd love to accept your patches! Before we can take them, we have to jump a couple of legal hurdles.

Please fill out either the individual or corporate Contributor License Agreement (CLA). More information about the CLA
and instructions for signing it [can be found here](https://github.com/kubernetes/community/blob/master/CLA.md).

***NOTE***: Only original source code from you and other people that have signed the CLA can be accepted into the
*repository.

## Finding Things That Need Help

If you're new to the project and want to help, but don't know where to start, we have a semi-curated list of issues that
should not need deep knowledge of the system. [Have a look and see if anything sounds
interesting](https://github.com/kubernetes-sigs/cluster-api/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22).
Alternatively, read some of the docs on other controllers and try to write your own, file and fix any/all issues that
come up, including gaps in documentation!

## Contributing a Patch

1. If you haven't already done so, sign a Contributor License Agreement (see details above).
1. Fork the desired repo, develop and test your code changes.
1. Submit a pull request.
    1. All code PR must be labeled with one of
        - ⚠️ (:warning:, major or breaking changes)
        - ✨ (:sparkles:, minor or feature additions)
        - 🐛 (:bug:, patch and bugfixes)
        - 📖 (:book:, documentation or proposals)
        - 🏃 (:running:, other)

All changes must be code reviewed. Coding conventions and standards are explained in the official [developer
docs](https://github.com/kubernetes/community/tree/master/contributors/devel). Expect reviewers to request that you
avoid common [go style mistakes](https://github.com/golang/go/wiki/CodeReviewComments) in your PRs.

## Backporting a Patch

Cluster API maintains older versions through `release-X.Y` branches. We accept backports of bug fixes to the most recent
release branch. For example, if the most recent branch is `release-0.2`, and the `master` branch is under active
development for v0.3.0, a bug fix that merged to `master` that also affects `v0.2.x` may be considered for backporting
to `release-0.2`. We generally do not accept PRs against older release branches.

## Features and bugs

Open [issues](https://github.com/kubernetes-sigs/cluster-api/issues/new/choose) to report bugs, or minor features.

For big feature, API and contract amendments, we follow the CAEP process as outlined below.

## Proposal process (CAEP)

The Cluster API Enhacement Proposal is the process this project uses to adopt new features, or changes to the APIs.

- The template, and accepted proposals live under `docs/proposals`.
- A proposal SHOULD be introduced and discussed during the weekly community meetings,
  [Kubernetes SIG Cluster Lifecycle mailing list](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle),
  or [discuss forum](https://discuss.kubernetes.io/c/contributors/cluster-api/).
- A proposal SHOULD be submitted first to the community using a collaborative writing platform, preferably Google Docs.
  - When using Google Docs, share the document with edit permissions for the [Kubernetes SIG Cluster Lifecycle mailing list](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle).

## Experiments

Proof of concepts, code experiments, or other initiatives can live under the `exp` folder and behind a feature gate.

- Experiments SHOULD not modify any of the publicly exposed APIs (e.g. CRDs).
- Experiments SHOULD not modify any existing CRD types outside of the experimental API group(s).
- Experiments SHOULD not modify any existing command line contracts.
- Experiments MUST not cause any breaking changes to existing (non-experimental) Go APIs.
- Experiments SHOULD introduce utility helpers in the go APIs for experiments that cross multiple components
  and require support from bootstrap, control plane, or infrastructure providers.
- Experiments follow a strict lifecycle: Alpha -> Beta prior to Graduation.
  - Alpha-stage experiments:
    - SHOULD not be enabled by default and any feature gates MUST be marked as 'Alpha' 
    - MUST be associated with a CAEP that is merged and in at least a provisional state
    - MAY be considered inactive and marked as deprecated if the following does not happen within the course of 1 minor release cycle:
      - Transition to Beta-stage
      - Active development towards progressing to Beta-stage
      - Either direct or downstream user evaluation
    - Any deprecated Alpha-stage experiment MAY be removed in the next minor release. 
  - Beta-stage experiments:
    - SHOULD be enabled by default, and any feature gates MUST be marked as 'Beta'
    - MUST be associated with a CAEP that is at least in the experimental state
    - MUST support conversions for any type changes
    - MUST remain backwards compatible unless updates are coinciding with a breaking Cluster API release
    - MAY be considered inactive and marked as deprecated if the following does not happen within the course of 1 minor release cycle:
      - Graduate
      - Active development towards Graduation
      - Either direct or downstream user consumption
    - Any deprecated Beta-stage experiment MAY be removed after being deprecated for an entire minor release. 
- Experiment Graduation MUST coincide with a breaking Cluster API release
- Experiment Graduation checklist:
  - [ ] MAY provide a way to be disabled, any feature gates MUST be marked as 'GA'
  - [ ] MUST undergo a full Kubernetes-style API review and update the CAEP with the plan to address any issues raised
  - [ ] CAEP MUST be in an implementable state and is fully up to date with the current implementation
  - [ ] CAEP MUST define transition plan for moving out of the experimental api group and code directories
  - [ ] CAEP MUST define any upgrade steps required for Existing Management and Workload Clusters
  - [ ] CAEP MUST define any upgrade steps required to be implemented by out-of-tree bootstrap, control plane, and infrastructure providers.

## Breaking Changes

Breaking changes are generally allowed in the `master` branch, as this is the branch used to develop the next minor
release of Cluster API.

There may be times, however, when `master` is closed for breaking changes. This is likely to happen as we near the
release of a new minor version.

Breaking changes are not allowed in release branches, as these represent minor versions that have already been released.
These versions have consumers who expect the APIs, behaviors, etc. to remain stable during the life time of the patch
stream for the minor release.

Examples of breaking changes include:

- Removing or renaming a field in a CRD
- Removing or renaming a CRD
- Removing or renaming an exported constant, variable, type, or function
- Updating the version of critical libraries such as controller-runtime, client-go, apimachinery, etc.
    - Some version updates may be acceptable, for picking up bug fixes, but maintainers must exercise caution when
      reviewing.

There may, at times, need to be exceptions where breaking changes are allowed in release branches. These are at the
discretion of the project's maintainers, and must be carefully considered before merging. An example of an allowed
breaking change might be a fix for a behavioral bug that was released in an initial minor version (such as `v0.3.0`).

## Merge Approval

Please see the [Kubernetes community document on pull
requests](https://git.k8s.io/community/contributors/guide/pull-requests.md) for more information about the merge
process.

## Google Doc Viewing Permissions

To gain viewing permissions to google docs in this project, please join either the
[kubernetes-dev](https://groups.google.com/forum/#!forum/kubernetes-dev) or
[kubernetes-sig-cluster-lifecycle](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle) google
group.

## Issue and Pull Request Management

Anyone may comment on issues and submit reviews for pull requests. However, in order to be assigned an issue or pull
request, you must be a member of the [Kubernetes SIGs](https://github.com/kubernetes-sigs) GitHub organization.

If you are a Kubernetes GitHub organization member, you are eligible for membership in the Kubernetes SIGs GitHub
organization and can request membership by [opening an
issue](https://github.com/kubernetes/org/issues/new?template=membership.md&title=REQUEST%3A%20New%20membership%20for%20%3Cyour-GH-handle%3E)
against the kubernetes/org repo.

However, if you are a member of any of the related Kubernetes GitHub organizations but not of the Kubernetes org, you
will need explicit sponsorship for your membership request. You can read more about Kubernetes membership and
sponsorship [here](https://github.com/kubernetes/community/blob/master/community-membership.md).

Cluster API maintainers can assign you an issue or pull request by leaving a `/assign <your Github ID>` comment on the
issue or pull request.
