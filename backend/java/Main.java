package server;

public class Main{
	public static void main(String[] args) throws Exception{
		SpiderServer ss = new SpiderServer(9999);
		ss.listen();
	}
}
