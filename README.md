# Whip

The `whip` tool was born with the mission to simplify the integration of the actions to
follow after the installation of a package and/or to manage OS checks.

It uses YAML to handle the triggers/hooks to run.

In the Gentoo/Funtoo worlds this means for the post-install phase to execute the commands
now available inside the ebuilds in the `postinst` section.

## Specifications

Hereinafter, an example of the YAML content of the actions to run:

```yaml

hooks:
  elogind_postinst:
    remediate: elogind_setup
    check: elogind_valid
    description: | 
      blah blah blah. This is what it does. This is why it is needed.
    keywords:
      - core
      - elogind
    actions:
      # Action executed in &&.
      - do something
  elogind_postrm:
    remediate: elogind_postrm
    check: elogind_cleanup
    description: | 
      blah blah blah. This is what it does. This is why it is needed.
    actions:
      # Action executed in &&.
      - do something
    depends:
      - hooks_1
    keywords:
      - core
      - elogind

  mime_update:
    remediate: x_setup
    check: xorg_valid
    description: |
      Update mime and gdk database/cache.
    actions:
      - source /etc/profile && update-mime-database /usr/share/mime
      - source /etc/profile && gdk-pixbuf-query-loaders --update-cache
    keywords:
      - x

  gtk_update:
    remediate: x_setup
    check: xorg_valid
    description: |
      Update gdk and glib database/schemas.
    actions:
      # Fix gnome icons caches
      - rm -f /usr/share/icons/hicolor/icon-theme.cache
      - source /etc/profile && gtk-update-icon-cache -f /usr/share/icons/*
      - source /etc/profile && glib-compile-schemas /usr/share/glib-2.0/schemas
    keywords:
      - x

  fonts_create_scale:
    remediate: x_setup
    check: xorg_valid
    description: |
      Create fonts.scale file
    envs:
      # Required envs
      - FONT_DIR
    actions:
      - >-
        if [[ ${FONT_DIR} != Speedo && ${FONT_DIR} != CID ]]; then
        echo "Generating fonts.scale for ${FONT_DIR}" ;
        mkfontscale
        -a "${EROOT}/usr/share/fonts/encodings/encodings.dir"
        -- "${EROOT}/usr/share/fonts/${FONT_DIR}" ;
        fi

    keywords:
      - x
      - fonts

  fonts_create_dir:
    remediate: x_setup
    check: xorg_valid
    description: |
      Create fonts dir.
    envs:
      # Required envs
      - FONT_DIR
    actions:
      - >-
        echo "Generating fonts.dir for ${FONT_DIR}"
        mkfontdir
        -e "${EROOT}"/usr/share/fonts/encodings
        -e "${EROOT}"/usr/share/fonts/encodings/large
        -- "${EROOT}/usr/share/fonts/${FONT_DIR}"
    keywords:
      - x
      - fonts

  fonts_setup:
    remediate: x_setup
    check: xorg_valid
    description: |
      Setup all available fonts
    actions:
      - >-
        for i in $(ls ${EROOT}/usr/share/fonts/) ; do
        if [ $i != "encoding" ] ; then
        FONT_DIR=$i whip hook fonts_create_scale ;
        FONT_DIR=$i whip hook fonts_create_dir ;
        fi ;
        done
    keywords:
      - x
      - fonts

```

## Commands

```bash
# will call all hooks with remediate x_setup
$> whip remediate x_setup

# Will call a specific hook
$> whip hook elogind.elogind_postrm
```
