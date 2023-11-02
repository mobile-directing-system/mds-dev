# Configuration file for the Sphinx documentation builder.
#
# This file only contains a selection of the most common options. For a full
# list see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Path setup --------------------------------------------------------------

# If extensions (or modules to document with autodoc) are in another directory,
# add these directories to sys.path here. If the directory is relative to the
# documentation root, use os.path.abspath to make it absolute, like shown here.
#
# import os
# import sys
# sys.path.insert(0, os.path.abspath('.'))


# -- Project information -----------------------------------------------------

project = 'MDS Server'
copyright = '2023, Yves Haas, Laurin Todt, Lennart Altenhof'
author = 'Yves Haas, Laurin Todt, Lennart Altenhof'


# -- General configuration ---------------------------------------------------

# Add any Sphinx extension module names here, as strings. They can be
# extensions coming with Sphinx (named 'sphinx.ext.*') or your custom
# ones.
extensions = [
    "sphinx.ext.duration",
    "sphinx.ext.extlinks",
    'sphinxemoji.sphinxemoji',
]

# Add external URLs.

extlinks = {
    "docker-homepage": ("https://www.docker.com/%s", "Docker"),
    "docker-install": ("https://docs.docker.com/get-docker/%s", "Install Docker"),
    "e-164": ("https://en.wikipedia.org/wiki/E.164%s", "E.164"),
    "fluent-bit-homepage": ("https://fluentbit.io/%s", "Fluent Bit"),
    "git-homepage": ("https://git-scm.com/%s", "Git"),
    "github-repo": ("https://github.com/mobile-directing-system/mds-server%s", "GitHub Repository"),
    "goland-homepage": ("https://www.jetbrains.com/go/%s", "GoLand"),
    "intellij-cloud-code-plugin-homepage": ("https://plugins.jetbrains.com/plugin/8079-cloud-code%s", "IntelliJ Cloud Code Plugin"),
    "intellij-cloud-code-plugin-install": ("https://cloud.google.com/code/docs/intellij/install%s", "Cloud Code Instructions"),
    "minikube-homepage": ("https://minikube.sigs.k8s.io/%s", "minikube Homepage"),
    "minikube-install": ("https://minikube.sigs.k8s.io/docs/start/%s", "Install minikube"),
    "skaffold-homepage": ("https://skaffold.dev/%s", "Skaffold Homepage"),
    "skaffold-install": ("https://skaffold.dev/docs/install/%s", "Install Skaffold"),
}

# Add any paths that contain templates here, relative to this directory.
templates_path = ['_templates']

# List of patterns, relative to source directory, that match files and
# directories to ignore when looking for source files.
# This pattern also affects html_static_path and html_extra_path.
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']


# -- Options for HTML output -------------------------------------------------

# The theme to use for HTML and HTML Help pages.  See the documentation for
# a list of builtin themes.
#
html_theme = 'furo'

# Add any paths that contain custom static files (such as style sheets) here,
# relative to this directory. They are copied after the builtin static files,
# so a file named "default.css" will overwrite the builtin "default.css".
html_static_path = ['_static']
