import { setItem, getItem, removeItem } from '../localStorage';

it('should set the value in key foo with prefix', () => {
  vi.spyOn(Object.getPrototypeOf(window.localStorage), 'setItem');

  setItem('foo', 'bar');
  expect(localStorage.setItem).toBeCalledWith('dash-ops:foo', 'bar');
});

it('should get the value in key foo', () => {
  vi.spyOn(Object.getPrototypeOf(window.localStorage), 'getItem');

  const value = getItem('foo');
  expect(localStorage.getItem).toBeCalledWith('dash-ops:foo');
  expect(value).toEqual('bar');
});

it('should remove the value in key foo', () => {
  vi.spyOn(Object.getPrototypeOf(window.localStorage), 'removeItem');

  removeItem('foo');
  expect(localStorage.removeItem).toBeCalledWith('dash-ops:foo');
});
