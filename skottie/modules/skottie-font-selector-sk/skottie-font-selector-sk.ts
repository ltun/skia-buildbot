/**
 * @module skottie-font-selector-sk
 * @description <h2><code>skottie-font-selector-sk</code></h2>
 *
 * <p>
 *   A font selector to modify text layers from a skottie animation.
 *   The list of fonts (availableFonts) has been selected
 *   from a larger set of fonts available in a mirror of google web fonts.
 *   Refer here for more information:
 *   https://skia.googlesource.com/buildbot/+/refs/heads/main/skottie/modules/skottie-sk/skottie-sk.ts#1060
 * </p>
 *
 */
import { define } from '../../../elements-sk/modules/define';
import { html, TemplateResult } from 'lit-html';
import { ElementSk } from '../../../infra-sk/modules/ElementSk';
import '../skottie-dropdown-sk';
import { DropdownSelectEvent } from '../skottie-dropdown-sk/skottie-dropdown-sk';
import { LottieAnimation, LottieLayer, LottieAsset } from '../types';

type FontType = {
  fName: string;
  fStyle: string;
  fFamily: string;
};

type OptionType = {
  value: string;
  id: string;
  selected?: boolean;
};

const LAYER_TYPE_TEXT = 5;

const availableFonts: FontType[] = [
  {
    fName: 'Righteous-Regular',
    fStyle: 'Righteous',
    fFamily: 'Regular',
  },
  {
    fName: 'BarlowCondensed-Regular',
    fStyle: 'Regular',
    fFamily: 'BarlowCondensed',
  },
  {
    fName: 'Anton-Regular',
    fStyle: 'Regular',
    fFamily: 'Anton',
  },
  {
    fName: 'DMSans-Regular',
    fStyle: 'Regular',
    fFamily: 'DMSans',
  },
  {
    fName: 'KronaOne-Regular',
    fStyle: 'Regular',
    fFamily: 'KronaOne',
  },
  {
    fName: 'BarlowCondensed-SemiBold',
    fStyle: 'SemiBold',
    fFamily: 'BarlowCondensed',
  },
  {
    fName: 'Archivo-BoldItalic',
    fStyle: 'BoldItalic',
    fFamily: 'Archivo',
  },
  {
    fName: 'Montserrat-Bold',
    fStyle: 'Bold',
    fFamily: 'Montserrat',
  },
  {
    fName: 'Syncopate-Bold',
    fStyle: 'Bold',
    fFamily: 'Syncopate',
  },
  {
    fName: 'SairaCondensed-ExtraBold',
    fStyle: 'ExtraBold',
    fFamily: 'SairaCondensed',
  },
  {
    fName: 'LuckiestGuy-Regular',
    fStyle: 'Regular',
    fFamily: 'LuckiestGuy',
  },
  {
    fName: 'Tomorrow-ExtraBoldItalic',
    fStyle: 'ExtraBoldItalic',
    fFamily: 'Tomorrow',
  },
  {
    fName: 'LondrinaSolid-Black',
    fStyle: 'Black',
    fFamily: 'LondrinaSolid',
  },
  {
    fName: 'Montserrat-Black',
    fStyle: 'Black',
    fFamily: 'Montserrat',
  },
  {
    fName: 'TitilliumWeb-Black',
    fStyle: 'Black',
    fFamily: 'TitilliumWeb',
  },
  {
    fName: 'Poppins-BlackItalic',
    fStyle: 'BlackItalic',
    fFamily: 'Poppins',
  },
  {
    fName: 'Comfortaa-Light',
    fStyle: 'Light',
    fFamily: 'Comfortaa',
  },
  {
    fName: 'Boogaloo-Regular',
    fStyle: 'Regular',
    fFamily: 'Boogaloo',
  },
  {
    fName: 'Chewy-Regular',
    fStyle: 'Regular',
    fFamily: 'Chewy',
  },
  {
    fName: 'Overlock-BlackItalic',
    fStyle: 'BlackItalic',
    fFamily: 'Overlock',
  },
  {
    fName: 'FredokaOne-Regular',
    fStyle: 'Regular',
    fFamily: 'FredokaOne',
  },
  {
    fName: 'Shrikhand-Regular',
    fStyle: 'Regular',
    fFamily: 'Shrikhand',
  },
  {
    fName: 'SpicyRice-Regular',
    fStyle: 'Regular',
    fFamily: 'SpicyRice',
  },
  {
    fName: 'Modak-Regular',
    fStyle: 'Regular',
    fFamily: 'Modak',
  },
  {
    fName: 'Chango-Regular',
    fStyle: 'Regular',
    fFamily: 'Chango',
  },
  {
    fName: 'Sniglet-ExtraBold',
    fStyle: 'ExtraBold',
    fFamily: 'Sniglet',
  },
  {
    fName: 'AmaticSC-Bold',
    fStyle: 'Bold',
    fFamily: 'AmaticSC',
  },
  {
    fName: 'CaveatBrush-Regular',
    fStyle: 'Regular',
    fFamily: 'CaveatBrush',
  },
  {
    fName: 'CoveredByYourGrace-Regular',
    fStyle: 'Regular',
    fFamily: 'CoveredByYourGrace',
  },
  {
    fName: 'Knewave-Regular',
    fStyle: 'Regular',
    fFamily: 'Knewave',
  },
  {
    fName: 'PermanentMarker-Regular',
    fStyle: 'Regular',
    fFamily: 'PermanentMarker',
  },
  {
    fName: 'Damion-Regular',
    fStyle: 'Regular',
    fFamily: 'Damion',
  },
  {
    fName: 'Neonderthaw-Regular',
    fStyle: 'Regular',
    fFamily: 'Neonderthaw',
  },
  {
    fName: 'Pacifico-Regular',
    fStyle: 'Regular',
    fFamily: 'Pacifico',
  },
  {
    fName: 'Lobster-Regular',
    fStyle: 'Regular',
    fFamily: 'Lobster',
  },
  {
    fName: 'Molle-Regular',
    fStyle: 'Regular',
    fFamily: 'Molle',
  },
  {
    fName: 'Bahiana-Regular',
    fStyle: 'Regular',
    fFamily: 'Bahiana',
  },
  {
    fName: 'JollyLodger-Regular',
    fStyle: 'Regular',
    fFamily: 'JollyLodger',
  },
  {
    fName: 'LifeSavers-ExtraBold',
    fStyle: 'ExtraBold',
    fFamily: 'LifeSavers',
  },
  {
    fName: 'Warnes-Regular',
    fStyle: 'Regular',
    fFamily: 'Warnes',
  },
  {
    fName: 'Ranchers-Regular',
    fStyle: 'Regular',
    fFamily: 'Ranchers',
  },
  {
    fName: 'Creepster-Regular',
    fStyle: 'Regular',
    fFamily: 'Creepster',
  },
  {
    fName: 'Slackey-Regular',
    fStyle: 'Regular',
    fFamily: 'Slackey',
  },
  {
    fName: 'Monoton-Regular',
    fStyle: 'Regular',
    fFamily: 'Monoton',
  },
  {
    fName: 'NewRocker-Regular',
    fStyle: 'Regular',
    fFamily: 'NewRocker',
  },
  {
    fName: 'ChelaOne-Regular',
    fStyle: 'Regular',
    fFamily: 'ChelaOne',
  },
  {
    fName: 'GermaniaOne-Regular',
    fStyle: 'Regular',
    fFamily: 'GermaniaOne',
  },
  {
    fName: 'Metamorphous-Regular',
    fStyle: 'Regular',
    fFamily: 'Metamorphous',
  },
  {
    fName: 'Spectral-BoldItalic',
    fStyle: 'BoldItalic',
    fFamily: 'Spectral',
  },
  {
    fName: 'Corben-Bold',
    fStyle: 'Bold',
    fFamily: 'Corben',
  },
  {
    fName: 'EBGaramond-Medium',
    fStyle: 'Medium',
    fFamily: 'EBGaramond',
  },
  {
    fName: 'PlayfairDisplay-SemiBoldItalic',
    fStyle: 'SemiBoldItalic',
    fFamily: 'PlayfairDisplay',
  },
  {
    fName: 'Merriweather-Regular',
    fStyle: 'Regular',
    fFamily: 'Merriweather',
  },
  {
    fName: 'AbrilFatface-Regular',
    fStyle: 'Regular',
    fFamily: 'AbrilFatface',
  },
  {
    fName: 'TextMeOne-Regular',
    fStyle: 'Regular',
    fFamily: 'TextMeOne',
  },
];

export interface SkottieFontEventDetail {
  animation: LottieAnimation;
}

export class SkottieFontSelectorSk extends ElementSk {
  private _animation: LottieAnimation | null = null;

  private static template = (
    ele: SkottieFontSelectorSk
  ): TemplateResult => html`
    <div class="wrapper">
      <skottie-dropdown-sk
        id="view-exporter"
        .name="dropdown-exporter"
        .options=${ele.buildFontOptions()}
        @select=${ele.fontTypeSelectHandler}
        border
        full
      >
      </skottie-dropdown-sk>
    </div>
  `;

  constructor() {
    super(SkottieFontSelectorSk.template);
  }

  connectedCallback(): void {
    super.connectedCallback();
    this._render();
  }

  disconnectedCallback(): void {
    super.disconnectedCallback();
  }

  private buildFontOptions(): OptionType[] {
    const fontOptions: OptionType[] = availableFonts.map((font) => {
      const isSelected =
        this._animation?.fonts?.list &&
        this._animation?.fonts?.list[0].fName === font.fName;
      return { id: font.fName, value: font.fName, selected: isSelected };
    });
    fontOptions.unshift({
      id: '',
      value: 'Select Font',
    });
    return fontOptions;
  }

  private updateFontInLayers(
    layers: LottieLayer[],
    targetFont: string,
    replacingFont: string
  ): void {
    layers.forEach((layer) => {
      if (layer.ty === LAYER_TYPE_TEXT) {
        if (layer.t?.d.k[0].s.f === targetFont) {
          layer.t.d.k[0].s.f = replacingFont;
        }
      }
    });
  }

  private updateFontInAssets(
    assets: LottieAsset[],
    targetFont: string,
    replacingFont: string
  ): void {
    assets.forEach((asset) => {
      if (asset.layers) {
        this.updateFontInLayers(asset.layers, targetFont, replacingFont);
      }
    });
  }

  private fontTypeSelectHandler(ev: CustomEvent<DropdownSelectEvent>): void {
    // This event handler replaces the animation in place
    // instead of creating a copy of the lottie animation.
    // If there is a reason why it should create a copy, this can be updated
    if (ev.detail.value && this._animation) {
      if (this._animation.fonts?.list?.length) {
        // We're changing a single font for the time being.
        // If we need to change multiple fonts, we probably will want to add a dropdown per font.
        const currentFontName = this._animation.fonts.list[0].fName;
        const newFontData = availableFonts.find(
          (font) => font.fName === ev.detail.value
        );
        if (newFontData) {
          this._animation.fonts.list[0].fName = newFontData.fName;
          this._animation.fonts.list[0].fFamily = newFontData.fFamily;
          this._animation.fonts.list[0].fStyle = newFontData.fStyle;
          this.updateFontInLayers(
            this._animation.layers,
            currentFontName,
            newFontData.fName
          );
          this.updateFontInAssets(
            this._animation.assets,
            currentFontName,
            newFontData.fName
          );
        }

        this.dispatchEvent(
          new CustomEvent<SkottieFontEventDetail>('animation-updated', {
            detail: {
              animation: this._animation,
            },
            bubbles: true,
          })
        );
      }
    }
  }

  set animation(value: LottieAnimation) {
    this._animation = value;
    this._render();
  }
}

define('skottie-font-selector-sk', SkottieFontSelectorSk);