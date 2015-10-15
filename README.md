# go-gae-redirect
A simple App Engine application to redirect between http://pawelsz.eu/g -> http://github.com/orian

I use it to redirect paths on my domain to my profiles on other websites.

The deployment is super easy, just create a Google App Engine instance, copy `app.yaml.orig` to `app.yaml` and specify application name. Then you need to deploy with appengine deployment scripts.
