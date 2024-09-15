/* eslint-disable no-template-curly-in-string */
import {
  LocalizationKey, LocalizationBodyKey, LocalizationTitleKey, LanguageCode,
} from './types';

export const LOCALIZED_MESSAGES: Record<LanguageCode, Record<LocalizationKey, string>> = {
  en: {
    [LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE]: 'Deposit Successful',
    [LocalizationBodyKey.DEPOSIT_SUCCESS_BODY]: 'You have successfully deposited {AMOUNT} {MARKET} to your dYdX account.',
    [LocalizationTitleKey.ORDER_FILLED_TITLE]: 'Order Filled',
    [LocalizationBodyKey.ORDER_FILLED_BODY]: 'Your order for {AMOUNT} {MARKET} was filled at ${AVERAGE_PRICE}',
    [LocalizationTitleKey.ORDER_TRIGGERED_TITLE]: '{MARKET} Order Triggered',
    [LocalizationBodyKey.ORDER_TRIGGERED_BODY]: 'Your order for {AMOUNT} {MARKET} was triggered at ${PRICE}',
  },
  es: {
    [LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE]: 'Depósito Exitoso',
    [LocalizationBodyKey.DEPOSIT_SUCCESS_BODY]: 'Has depositado exitosamente {AMOUNT} {MARKET} en tu cuenta dYdX.',
    [LocalizationTitleKey.ORDER_FILLED_TITLE]: 'Orden Ejecutada',
    [LocalizationBodyKey.ORDER_FILLED_BODY]: 'Tu orden de {AMOUNT} {MARKET} se ejecutó a ${AVERAGE_PRICE}',
    [LocalizationTitleKey.ORDER_TRIGGERED_TITLE]: 'Orden de {MARKET} Activada',
    [LocalizationBodyKey.ORDER_TRIGGERED_BODY]: 'Tu orden de {AMOUNT} {MARKET} se activó a ${PRICE}',
  },
  fr: {
    [LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE]: 'Dépôt Réussi',
    [LocalizationBodyKey.DEPOSIT_SUCCESS_BODY]: 'Vous avez déposé avec succès {AMOUNT} {MARKET} sur votre compte dYdX.',
    [LocalizationTitleKey.ORDER_FILLED_TITLE]: 'Ordre Exécuté',
    [LocalizationBodyKey.ORDER_FILLED_BODY]: 'Votre ordre de {AMOUNT} {MARKET} a été exécuté à ${AVERAGE_PRICE}',
    [LocalizationTitleKey.ORDER_TRIGGERED_TITLE]: 'Ordre {MARKET} Déclenché',
    [LocalizationBodyKey.ORDER_TRIGGERED_BODY]: 'Votre ordre de {AMOUNT} {MARKET} a été déclenché à ${PRICE}',
  },
  de: {
    [LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE]: 'Einzahlung Erfolgreich',
    [LocalizationBodyKey.DEPOSIT_SUCCESS_BODY]: 'Sie haben erfolgreich {AMOUNT} {MARKET} auf Ihr dYdX-Konto eingezahlt.',
    [LocalizationTitleKey.ORDER_FILLED_TITLE]: 'Auftrag Ausgeführt',
    [LocalizationBodyKey.ORDER_FILLED_BODY]: 'Ihr Auftrag für {AMOUNT} {MARKET} wurde zu ${AVERAGE_PRICE} ausgeführt',
    [LocalizationTitleKey.ORDER_TRIGGERED_TITLE]: '{MARKET} Auftrag Ausgelöst',
    [LocalizationBodyKey.ORDER_TRIGGERED_BODY]: 'Ihr Auftrag für {AMOUNT} {MARKET} wurde bei ${PRICE} ausgelöst',
  },
  it: {
    [LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE]: 'Deposito Riuscito',
    [LocalizationBodyKey.DEPOSIT_SUCCESS_BODY]: 'Hai depositato con successo {AMOUNT} {MARKET} sul tuo account dYdX.',
    [LocalizationTitleKey.ORDER_FILLED_TITLE]: 'Ordine Eseguito',
    [LocalizationBodyKey.ORDER_FILLED_BODY]: 'Il tuo ordine di {AMOUNT} {MARKET} è stato eseguito a ${AVERAGE_PRICE}',
    [LocalizationTitleKey.ORDER_TRIGGERED_TITLE]: 'Ordine {MARKET} Attivato',
    [LocalizationBodyKey.ORDER_TRIGGERED_BODY]: 'Il tuo ordine di {AMOUNT} {MARKET} è stato attivato a ${PRICE}',
  },
  ja: {
    [LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE]: '入金成功',
    [LocalizationBodyKey.DEPOSIT_SUCCESS_BODY]: 'dYdXアカウントに{AMOUNT} {MARKET}を正常に入金しました。',
    [LocalizationTitleKey.ORDER_FILLED_TITLE]: '注文約定',
    [LocalizationBodyKey.ORDER_FILLED_BODY]: '{AMOUNT} {MARKET}の注文が${AVERAGE_PRICE}で約定しました',
    [LocalizationTitleKey.ORDER_TRIGGERED_TITLE]: '{MARKET}注文トリガー',
    [LocalizationBodyKey.ORDER_TRIGGERED_BODY]: '{AMOUNT} {MARKET}の注文が${PRICE}でトリガーされました',
  },
  ko: {
    [LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE]: '입금 성공',
    [LocalizationBodyKey.DEPOSIT_SUCCESS_BODY]: 'dYdX 계정에 {AMOUNT} {MARKET}을(를) 성공적으로 입금했습니다.',
    [LocalizationTitleKey.ORDER_FILLED_TITLE]: '주문 체결',
    [LocalizationBodyKey.ORDER_FILLED_BODY]: '{AMOUNT} {MARKET} 주문이 ${AVERAGE_PRICE}에 체결되었습니다',
    [LocalizationTitleKey.ORDER_TRIGGERED_TITLE]: '{MARKET} 주문 트리거',
    [LocalizationBodyKey.ORDER_TRIGGERED_BODY]: '{AMOUNT} {MARKET} 주문이 ${PRICE}에서 트리거되었습니다',
  },
  zh: {
    [LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE]: '存款成功',
    [LocalizationBodyKey.DEPOSIT_SUCCESS_BODY]: '您已成功将 {AMOUNT} {MARKET} 存入您的 dYdX 账户。',
    [LocalizationTitleKey.ORDER_FILLED_TITLE]: '订单已成交',
    [LocalizationBodyKey.ORDER_FILLED_BODY]: '您的 {AMOUNT} {MARKET} 订单已以 ${AVERAGE_PRICE} 成交',
    [LocalizationTitleKey.ORDER_TRIGGERED_TITLE]: '{MARKET} 订单已触发',
    [LocalizationBodyKey.ORDER_TRIGGERED_BODY]: '您的 {AMOUNT} {MARKET} 订单已在 ${PRICE} 触发',
  },
};
