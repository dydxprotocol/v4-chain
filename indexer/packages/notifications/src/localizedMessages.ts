/* eslint-disable no-template-curly-in-string */
import { LocalizationKey, LanguageCode } from './types';

export const LOCALIZED_MESSAGES: Record<LanguageCode, Record<LocalizationKey, string>> = {
  en: {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: 'Deposit Successful',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'You have successfully deposited {AMOUNT} {MARKET} to your dYdX account.',
    [LocalizationKey.ORDER_FILLED_TITLE]: 'Order Filled',
    [LocalizationKey.ORDER_FILLED_BODY]: 'Your order for {AMOUNT} {MARKET} was filled at ${AVERAGE_PRICE}',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: '{MARKET} Order Triggered',
    [LocalizationKey.ORDER_TRIGGERED_BODY]: 'Your order for {AMOUNT} {MARKET} was triggered at ${PRICE}',
  },
  es: {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: 'Depósito Exitoso',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'Has depositado exitosamente {AMOUNT} {MARKET} en tu cuenta dYdX.',
    [LocalizationKey.ORDER_FILLED_TITLE]: 'Orden Ejecutada',
    [LocalizationKey.ORDER_FILLED_BODY]: 'Tu orden de {AMOUNT} {MARKET} se ejecutó a ${AVERAGE_PRICE}',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: 'Orden de {MARKET} Activada',
    [LocalizationKey.ORDER_TRIGGERED_BODY]: 'Tu orden de {AMOUNT} {MARKET} se activó a ${PRICE}',
  },
  fr: {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: 'Dépôt Réussi',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'Vous avez déposé avec succès {AMOUNT} {MARKET} sur votre compte dYdX.',
    [LocalizationKey.ORDER_FILLED_TITLE]: 'Ordre Exécuté',
    [LocalizationKey.ORDER_FILLED_BODY]: 'Votre ordre de {AMOUNT} {MARKET} a été exécuté à ${AVERAGE_PRICE}',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: 'Ordre {MARKET} Déclenché',
    [LocalizationKey.ORDER_TRIGGERED_BODY]: 'Votre ordre de {AMOUNT} {MARKET} a été déclenché à ${PRICE}',
  },
  de: {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: 'Einzahlung Erfolgreich',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'Sie haben erfolgreich {AMOUNT} {MARKET} auf Ihr dYdX-Konto eingezahlt.',
    [LocalizationKey.ORDER_FILLED_TITLE]: 'Auftrag Ausgeführt',
    [LocalizationKey.ORDER_FILLED_BODY]: 'Ihr Auftrag für {AMOUNT} {MARKET} wurde zu ${AVERAGE_PRICE} ausgeführt',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: '{MARKET} Auftrag Ausgelöst',
    [LocalizationKey.ORDER_TRIGGERED_BODY]: 'Ihr Auftrag für {AMOUNT} {MARKET} wurde bei ${PRICE} ausgelöst',
  },
  it: {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: 'Deposito Riuscito',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'Hai depositato con successo {AMOUNT} {MARKET} sul tuo account dYdX.',
    [LocalizationKey.ORDER_FILLED_TITLE]: 'Ordine Eseguito',
    [LocalizationKey.ORDER_FILLED_BODY]: 'Il tuo ordine di {AMOUNT} {MARKET} è stato eseguito a ${AVERAGE_PRICE}',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: 'Ordine {MARKET} Attivato',
    [LocalizationKey.ORDER_TRIGGERED_BODY]: 'Il tuo ordine di {AMOUNT} {MARKET} è stato attivato a ${PRICE}',
  },
  ja: {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: '入金成功',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'dYdXアカウントに{AMOUNT} {MARKET}を正常に入金しました。',
    [LocalizationKey.ORDER_FILLED_TITLE]: '注文約定',
    [LocalizationKey.ORDER_FILLED_BODY]: '{AMOUNT} {MARKET}の注文が${AVERAGE_PRICE}で約定しました',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: '{MARKET}注文トリガー',
    [LocalizationKey.ORDER_TRIGGERED_BODY]: '{AMOUNT} {MARKET}の注文が${PRICE}でトリガーされました',
  },
  ko: {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: '입금 성공',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'dYdX 계정에 {AMOUNT} {MARKET}을(를) 성공적으로 입금했습니다.',
    [LocalizationKey.ORDER_FILLED_TITLE]: '주문 체결',
    [LocalizationKey.ORDER_FILLED_BODY]: '{AMOUNT} {MARKET} 주문이 ${AVERAGE_PRICE}에 체결되었습니다',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: '{MARKET} 주문 트리거',
    [LocalizationKey.ORDER_TRIGGERED_BODY]: '{AMOUNT} {MARKET} 주문이 ${PRICE}에서 트리거되었습니다',
  },
  zh: {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: '存款成功',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: '您已成功将 {AMOUNT} {MARKET} 存入您的 dYdX 账户。',
    [LocalizationKey.ORDER_FILLED_TITLE]: '订单已成交',
    [LocalizationKey.ORDER_FILLED_BODY]: '您的 {AMOUNT} {MARKET} 订单已以 ${AVERAGE_PRICE} 成交',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: '{MARKET} 订单已触发',
    [LocalizationKey.ORDER_TRIGGERED_BODY]: '您的 {AMOUNT} {MARKET} 订单已在 ${PRICE} 触发',
  },
};
